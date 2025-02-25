package routes

import (
	"pedersandvoll/foosballapi/handlers"
	"pedersandvoll/foosballapi/middleware"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App, h *handlers.Handlers) {
	app.Post("/register", h.RegisterUser)
	app.Post("/login", h.LoginUser)

	api := app.Group("/api")
	api.Use(middleware.AuthRequired(h.JWTSecret))

	api.Post("/refresh", h.RefreshToken)
	api.Get("/users", h.GetUsers)

	api.Post("/org", h.CreateOrganization)
	api.Post("/join/org", h.JoinOrg)
	api.Post("/edit/org", h.EditOrgSettings)

	api.Post("/season", h.CreateSeason)

	api.Get("/lobbies", h.GetLobbies)
	api.Post("/lobby", h.CreateLobby)
}
