package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/tk-neng/demo-go-fiber/config"
	"github.com/tk-neng/demo-go-fiber/handler"
	"github.com/tk-neng/demo-go-fiber/middleware"
	"github.com/tk-neng/demo-go-fiber/utils"
)


func RouteInit(r *fiber.App) {
	r.Get("/api/public/*", static.New(config.ProjectRootPath+"/public"))
	r.Post("/api/login", handler.Login) 

	// route user
	r.Get("/api/user", handler.GetAllUser)
	r.Get("/api/user/:id", handler.GetUserByID)
	r.Post("/api/user", handler.CreateUser)


	// route provider
	r.Get("/api/provider", handler.GetAllProvider)


	// route country and category
	r.Get("/api/country", handler.GatAllCountry)
	r.Get("/api/category", handler.GetAllCategory)

	// route post
	r.Get("/api/announce", handler.GetAllAnnouncePost,middleware.Auth)
	r.Get("/api/subject", handler.GetAllPost)
	r.Get("/api/announce/:id", handler.GetAnnouncePostByID)
	r.Get("/api/subject/:id", handler.GetPostByID)
	r.Post("/api/announce", handler.CreateAnnouncePost,utils.HandleFileImage,utils.HandleFileAttach)
	r.Post("/api/subject", handler.CreatePost,utils.HandleFileImage)
	r.Put("/api/announce/:id", handler.UpdateAnnouncePost, utils.HandleFileImage, utils.HandleFileAttach)
	r.Put("/api/subject/:id", handler.UpdatePost, utils.HandleFileImage)
	r.Delete("/api/announce/:id", handler.DeleteAnnouncePost)
	r.Delete("/api/subject/:id", handler.DeletePost)
}