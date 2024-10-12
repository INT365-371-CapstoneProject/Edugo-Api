package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/handler"
)

func RouteInit(r *fiber.App) {
	r.Get("/users", handler.GetAllUser)
	r.Post("/users", handler.CreateUser)
}