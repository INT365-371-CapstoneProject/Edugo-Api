package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/tk-neng/demo-go-fiber/config"
	"github.com/tk-neng/demo-go-fiber/handler"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func RouteInit(r *fiber.App) {
	r.Get("/api/public/*", static.New(config.ProjectRootPath+"/public"))
	r.Get("/api/annouce", handler.GetAllAnnoucePost)
	// r.Get("/api/posts/:id", handler.GetPostByID)
	r.Post("/api/posts/create", handler.CreatePost,utils.HandleFileImage,utils.HandleFileAttach)
	// r.Put("/api/post/update/:id", handler.UpdatePost, utils.HandleFileImage, utils.HandleFileAttach)
	// r.Delete("/api/post/delete/:id", handler.DeletePost)
}