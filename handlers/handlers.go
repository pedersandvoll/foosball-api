package handlers

import (
	"database/sql"
	"log"
	"pedersandvoll/foosballapi/config"
	"pedersandvoll/foosballapi/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	db        *config.Database
	JWTSecret []byte
}

func NewHandlers(db *config.Database, jwtSecret string) *Handlers {
	return &Handlers{
		db:        db,
		JWTSecret: []byte(jwtSecret),
	}
}

func (h *Handlers) RefreshToken(c *fiber.Ctx) error {
	username := c.Locals("username").(string)

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(h.JWTSecret)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

type User struct {
	ID       int    `json:"userid"`
	UserName string `json:"username"`
}

func (h *Handlers) GetUsers(c *fiber.Ctx) error {
	rows, err := h.db.Query("SELECT userid, username FROM users")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database query failed"})
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.UserName,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan row"})
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error iterating over rows"})
	}

	return c.JSON(users)
}

type UserByName struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	UserId   string `json:"userid"`
}

func (h *Handlers) getUserByUsername(username string) (UserByName, error) {
	var password string
	var userid string
	query := "SELECT username, password, userid FROM users WHERE username=$1;"
	row := h.db.QueryRow(query, username)
	switch err := row.Scan(&username, &password, &userid); err {
	case sql.ErrNoRows:
		return UserByName{}, err
	case nil:
		return UserByName{UserName: username, Password: password, UserId: userid}, nil
	default:
		return UserByName{}, err
	}
}

type UserBody struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) RegisterUser(c *fiber.Ctx) error {
	var body UserBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.UserName == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING userid"
	var userID int
	err = h.db.QueryRow(query, body.UserName, hashedPassword).Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"userid":  userID,
	})
}

func (h *Handlers) LoginUser(c *fiber.Ctx) error {
	var body UserBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.UserName == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	userExist, err := h.getUserByUsername(body.UserName)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User or password are wrong",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	isValid := utils.VerifyPassword(body.Password, userExist.Password)
	if !isValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User or password are wrong",
		})
	}

	claims := jwt.MapClaims{
		"username": userExist.UserName,
		"userid":   userExist.UserId,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(h.JWTSecret)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

type NewOrg struct {
	Name      string `json:"name"`
	OrgSecret string `json:"orgsecret"`
}

func (h *Handlers) CreateOrganization(c *fiber.Ctx) error {
	var body NewOrg
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.Name == "" || body.OrgSecret == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and orgsecret is required",
		})
	}

	query := "INSERT INTO organizations (name, orgsecret, orgowner) VALUES ($1, $2, $3) RETURNING orgid"
	var orgID int

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["userid"].(string)

	err := h.db.QueryRow(query, body.Name, body.OrgSecret, userID).Scan(&orgID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create org",
		})
	}

	queryOrgSettings := "INSERT INTO organizationsettings (orgid, orgowner) VALUES ($1, $2)"
	_, err = h.db.Exec(queryOrgSettings, orgID, userID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create org settings",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Org created successfully",
		"orgid":   orgID,
	})
}

func (h *Handlers) getOrgBySecret(orgsecret string, c *fiber.Ctx) (string, error) {
	var orgID string

	query := "SELECT orgid FROM organizations WHERE orgsecret=$1"
	err := h.db.QueryRow(query, orgsecret).Scan(&orgID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No org with that secret",
		})
	}

	return orgID, nil
}

type JoinOrg struct {
	OrgSecret string `json:"orgsecret"`
}

func (h *Handlers) JoinOrg(c *fiber.Ctx) error {
	var body JoinOrg
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.OrgSecret == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Org secret is required",
		})
	}

	orgID, err := h.getOrgBySecret(body.OrgSecret, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get org",
		})
	}

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["userid"].(string)

	query := "UPDATE users SET activeorg = $1 WHERE userid = $2;"
	_, err = h.db.Exec(query, orgID, userID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create org settings",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Added user to organization",
	})
}
