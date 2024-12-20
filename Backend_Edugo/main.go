package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/route"
)

func main() {
	// Iinitial database
	database.DatabaseInit()

	app := fiber.New(fiber.Config{
		// 50MB
		BodyLimit: 50 * 1024 * 1024,
	})


	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	// Initial route
	route.RouteInit(app)
	port := os.Getenv("PORT")
	app.Listen(":" + port)
}
