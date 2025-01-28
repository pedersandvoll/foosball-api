package routes

import (
	"pedersandvoll/foosballapi/handlers"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App, h *handlers.Handlers) {
	app.Get("/users", h.GetUsers)
	app.Post("/register", h.RegisterUser)
	app.Post("/login", h.LoginUser)
}
