package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/handler"
)

func RouteInit(r *fiber.App) {
	r.Get("/api/posts", handler.GetAllPost)
	r.Post("/api/posts/create", handler.CreatePost)
}