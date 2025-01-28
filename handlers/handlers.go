package handlers

import (
	"database/sql"
	"github.com/golang-jwt/jwt/v5"
	"pedersandvoll/foosballapi/config"
	"pedersandvoll/foosballapi/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	db *config.Database
}

func NewHandlers(db *config.Database) *Handlers {
	return &Handlers{
		db: db,
	}
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
}

func (h *Handlers) getUserByUsername(username string) (UserByName, error) {
	var password string
	query := "SELECT username, password FROM users WHERE username=$1;"
	row := h.db.QueryRow(query, username)
	switch err := row.Scan(&username, &password); err {
	case sql.ErrNoRows:
		return UserByName{}, err
	case nil:
		return UserByName{UserName: username, Password: password}, nil
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
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
