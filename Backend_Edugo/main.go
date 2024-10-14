package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/route"
)

func main() {
	// Iinitial database
	database.DatabaseInit()

	app := fiber.New()

	// Initial route
	route.RouteInit(app)

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	port := os.Getenv("PORT")
	app.Listen(":" + port)
}
