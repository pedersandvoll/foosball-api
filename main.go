package main

import (
	"log"
	"pedersandvoll/foosballapi/config"
	"pedersandvoll/foosballapi/handlers"
	"pedersandvoll/foosballapi/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbConfig := config.NewConfig()
	db, err := config.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer db.Close()

	app := fiber.New()

	h := handlers.NewHandlers(db, dbConfig.JWTSecret)

	routes.Routes(app, h)

	app.Listen(":3000")
}
