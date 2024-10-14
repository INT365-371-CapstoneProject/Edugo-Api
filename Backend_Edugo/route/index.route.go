package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/handler"
)

func RouteInit(r *fiber.App) {
	r.Get("/users", handler.GetAllUser)
	r.Get("/api/posts", handler.GetAllPost)
	r.Get("/api/subjects", handler.GetAllSubject)
	r.Post("/users", handler.CreateUser)
}