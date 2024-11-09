package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/tk-neng/demo-go-fiber/config"
	"github.com/tk-neng/demo-go-fiber/handler"
)

func RouteInit(r *fiber.App) {
	r.Get("/public/*", static.New(config.ProjectRootPath+"/public"))
	r.Get("/api/posts", handler.GetAllPost)
	r.Get("/api/posts/:id", handler.GetPostByID)
	r.Post("/api/posts/create", handler.CreatePost)
	r.Put("/api/post/update/:id", handler.UpdatePost)
	r.Delete("/api/post/delete/:id", handler.DeletePost)
}