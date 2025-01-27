package handlers

import (
	"pedersandvoll/foosballapi/config"
	"pedersandvoll/foosballapi/utils"
	"strings"

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

func (h *Handlers) GetUsers(c *fiber.Ctx) error {
	type User struct {
		ID       int    `json:"userid"`
		UserName string `json:"username"`
	}

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

func (h *Handlers) getUserByUsername(username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	err := h.db.QueryRow(query, username).Scan(&exists)
	return exists, err
}

type RegisterBody struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) RegisterUser(c *fiber.Ctx) error {
	var body RegisterBody

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

	// Return success with the new user ID
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"userid":  userID,
	})
}
