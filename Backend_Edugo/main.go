package main

import (
	"os"
	"strings"

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
		AllowOriginsFunc: func(origin string) bool {
			return origin == "http://localhost:5173" ||
				origin == "https://capstone24.sit.kmutt.ac.th/un2" ||
				strings.HasPrefix(origin, "http://192.168.") || // อนุญาต local network
				strings.HasPrefix(origin, "http://10.0.2.") || // อนุญาต Android emulator
				strings.HasPrefix(origin, "http://localhost") // อนุญาต localhost ทุก port
		},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	// Initial route
	route.RouteInit(app)
	port := os.Getenv("PORT")
	app.Listen(":" + port)
}
