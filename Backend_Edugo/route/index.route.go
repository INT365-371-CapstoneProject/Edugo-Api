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
	// Public routes
	public := r.Group("/api")
	public.Get("/public/*", static.New(config.ProjectRootPath+"/public"))
	public.Post("/login", handler.Login)

	// forgot password
	public.Post("/auth/forgot-password", handler.ForgotPassword)
	public.Post("/auth/verify-otp", handler.VerifyOTP)

	// User routes
	userGroup := public.Group("/user")
	userGroup.Get("/", handler.GetAllUser)
	userGroup.Get("/:id", handler.GetUserByID)
	userGroup.Post("/", handler.CreateUser)

	// Provider routes
	providerGroup := public.Group("/provider")
	providerGroup.Get("/", handler.GetAllProvider)
	providerGroup.Post("/", handler.CreateProvider)

	// Metadata routes (country and category)
	metadataGroup := public.Group("")
	metadataGroup.Get("/country", handler.GatAllCountry)
	metadataGroup.Get("/category", handler.GetAllCategory)

	// Announcement routes
	announceGroup := public.Group("/announce", middleware.AuthProvider)
	announceGroup.Get("/", handler.GetAllAnnouncePost)
	announceGroup.Get("/:id", handler.GetAnnouncePostByID)
	announceGroup.Post("/", handler.CreateAnnouncePost, utils.HandleFileImage, utils.HandleFileAttach)
	announceGroup.Put("/:id", handler.UpdateAnnouncePost, utils.HandleFileImage, utils.HandleFileAttach)
	announceGroup.Delete("/:id", handler.DeleteAnnouncePost)

	// Subject routes
	subjectGroup := public.Group("/subject", middleware.AuthAny)
	subjectGroup.Get("/", handler.GetAllPost)
	subjectGroup.Get("/:id", handler.GetPostByID)
	subjectGroup.Post("/", handler.CreatePost, utils.HandleFileImage)
	subjectGroup.Put("/:id", handler.UpdatePost, utils.HandleFileImage)
	subjectGroup.Delete("/:id", handler.DeletePost)
}
