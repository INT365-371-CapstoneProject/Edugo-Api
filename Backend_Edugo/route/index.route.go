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
	r.Get("/api/annouce", handler.GetAllAnnouncePost)
	r.Get("/api/country", handler.GatAllCountry)
	r.Get("/api/category", handler.GetAllCategory)
	r.Get("/api/annouce/:id", handler.GetAnnouncePostByID)
	r.Post("/api/annouce/create", handler.CreateAnnouncePost,utils.HandleFileImage,utils.HandleFileAttach)
	r.Put("/api/annouce/update/:id", handler.UpdateAnnouncePost, utils.HandleFileImage, utils.HandleFileAttach)
	r.Delete("/api/annouce/delete/:id", handler.DeleteAnnouncePost)
}