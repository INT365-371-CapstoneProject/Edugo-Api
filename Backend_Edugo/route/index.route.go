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
	r.Get("/api/announce", handler.GetAllAnnouncePost)
	r.Get("/api/subject", handler.GetAllPost)
	r.Get("/api/country", handler.GatAllCountry)
	r.Get("/api/category", handler.GetAllCategory)
	r.Get("/api/announce/:id", handler.GetAnnouncePostByID)
	r.Get("/api/subject/:id", handler.GetPostByID)
	r.Post("/api/announce/add", handler.CreateAnnouncePost,utils.HandleFileImage,utils.HandleFileAttach)
	r.Post("/api/subject/add", handler.CreatePost,utils.HandleFileImage)
	r.Put("/api/announce/update/:id", handler.UpdateAnnouncePost, utils.HandleFileImage, utils.HandleFileAttach)
	r.Put("/api/subject/update/:id", handler.UpdatePost, utils.HandleFileImage)
	r.Delete("/api/announce/delete/:id", handler.DeleteAnnouncePost)
	r.Delete("/api/subject/delete/:id", handler.DeletePost)
}