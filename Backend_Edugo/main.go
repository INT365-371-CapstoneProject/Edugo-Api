package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/route"
)

func main() {
	// Iinitial database
	database.DatabaseInit()

	app := fiber.New()

	// Initial route
	route.RouteInit(app)

	app.Listen(":8080")
}
