package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/tk-neng/demo-go-fiber/config"
	"github.com/tk-neng/demo-go-fiber/handler"
	"github.com/tk-neng/demo-go-fiber/middleware"
	// "github.com/tk-neng/demo-go-fiber/utils"
)

func RouteInit(r *fiber.App) {
	// Public routes
	public := r.Group("/api")
	public.Get("/public/*", static.New(config.ProjectRootPath+"/public"))

	// Login routes
	public.Post("/login", handler.Login)

	// forgot password
	public.Post("/auth/forgot-password", handler.ForgotPassword)
	public.Post("/auth/verify-otp", handler.VerifyOTP)

	// User routes
	userGroup := public.Group("/user", middleware.PermissionCreate)
	// userGroup.Get("/", handler.GetAllUser)
	// userGroup.Get("/:id", handler.GetUserByID)
	userGroup.Post("/", handler.CreateUser)
	userGroup.Delete("/:id", handler.DeleteUser)

	// Provider routes
	providerGroup := public.Group("/provider", middleware.PermissionCreate)
	// providerGroup.Get("/", handler.GetAllProvider)
	// providerGroup.Get("/:id", handler.GetIDProvider)
	providerGroup.Post("/", handler.CreateProvider)

	// Metadata routes (country and category)
	metadataGroup := public.Group("", middleware.AuthAny)
	metadataGroup.Get("/country", handler.GatAllCountry)
	metadataGroup.Get("/category", handler.GetAllCategory)
	metadataGroup.Get("/provider/:id", handler.GetIDProvider)
	metadataGroup.Get("/provider/avatar/:id", handler.GetAvatarImageByProviderID)


	// Profile routes
	profileGroup := public.Group("/profile", middleware.AuthAny)
	profileGroup.Get("/", handler.GetProfile)
	profileGroup.Get("/avatar", handler.GetAvatarImage)
	profileGroup.Put("/", handler.UpdateProfile)
	profileGroup.Post("/change-password", handler.ChangePassword)

	// Search routes
	searchGroup := public.Group("/search", middleware.AuthAny)
	searchGroup.Get("/announce-provider", handler.SearchAnnouncementsForProvider)
	searchGroup.Get("/announce-admin", handler.SearchAnnouncementsForAdmin)
	searchGroup.Get("/announce-user", handler.SearchAnnouncementsForUser)
	searchGroup.Get("/subject", handler.SearchPosts)

	// Announcement for user routes
	announceUserGroup := public.Group("/announce-user", middleware.AuthAny)
	announceUserGroup.Get("/", handler.GetAllAnnouncePostForUser)
	announceUserGroup.Get("/:id", handler.GetAnnouncePostByIDForUser)
	announceUserGroup.Get("/:id/image", handler.GetAnnounceImage)
	announceUserGroup.Post("/bookmark", handler.GetAnnouncePostByBookmark)
	// Announcement for Admin routes
	announceAdminGroup := public.Group("/announce-admin", middleware.AuthAdmin)
	announceAdminGroup.Get("/", handler.GetAllAnnouncePostForAdmin)
	announceAdminGroup.Get("/:id", handler.GetAnnouncePostByIDForAdmin)
	announceAdminGroup.Get("/:id/image", handler.GetAnnounceImage)
	announceAdminGroup.Get("/:id/attach", handler.GetAnnouncePostAttach)
	announceAdminGroup.Delete("/:id", handler.DeleteAnnouncePostForAdmin)

	// Subject for Admin routes
	subjectAdminGroup := public.Group("/subject-admin", middleware.AuthAdmin)
	subjectAdminGroup.Delete("/:id", handler.DeletePostForAdmin)

	// Create Admin for SuperAdmin
	adminGroup := public.Group("/superadmin", middleware.AuthSuperAdmin)
	adminGroup.Post("/", handler.CreateAdminForSuperadmin)

	// Admin routes
	adminManageGroup := public.Group("/admin", middleware.AuthAdmin)
	adminManageGroup.Get("/user", handler.GetAllUser)
	adminManageGroup.Get("/user/:id", handler.GetIDUser)
	adminManageGroup.Get("/provider", handler.GetAllProviderForAdmin)
	adminManageGroup.Get("/provider/:id", handler.GetIDProviderForAdmin)
	adminManageGroup.Put("/verify/:id", handler.VerifyProviderForAdmin)
	adminManageGroup.Post("/manage-user", handler.ManageAllUser)
	adminManageGroup.Post("/edit", handler.EditAllUserForAdmin)

	// Announcement routes
	announceGroup := public.Group("/announce", middleware.AuthProvider)
	announceGroup.Get("/", handler.GetAllAnnouncePostForProvider)
	announceGroup.Get("/:id", handler.GetAnnouncePostByIDForProvider)
	announceGroup.Get("/:id/image", handler.GetAnnounceImage)
	announceGroup.Get("/:id/attach", handler.GetAnnouncePostAttach)
	announceGroup.Post("/", handler.CreateAnnouncePostForProvider)
	announceGroup.Put("/:id", handler.UpdateAnnouncePostForProvider)
	announceGroup.Delete("/:id", handler.DeleteAnnouncePostForProvider)

	// Subject routes
	subjectGroup := public.Group("/subject", middleware.AuthAny)
	subjectGroup.Get("/", handler.GetAllPost)
	subjectGroup.Get("/:id", handler.GetPostByID)
	subjectGroup.Get("/:id/image", handler.GetPostImage)
	subjectGroup.Get("/:id/avatar", handler.GetPostAvatarByAccountID)
	subjectGroup.Post("/", handler.CreatePost)
	subjectGroup.Put("/:id", handler.UpdatePost)
	subjectGroup.Delete("/:id", handler.DeletePost)

	// Comment routes
	commentGroup := public.Group("/comment", middleware.AuthAny)
	commentGroup.Get("/", handler.GetAllComment)
	commentGroup.Get("/post/:post_id", handler.GetCommentByPostID)
	commentGroup.Get("/:id/image", handler.GetCommentImage)
	commentGroup.Get("/:id/avatar", handler.GetCommentAvatarImageByAccountID)
	commentGroup.Post("/", handler.CreateComment)
	commentGroup.Delete("/:id", handler.DeleteComment)
	commentGroup.Put("/:id", handler.UpdateComment)

	// Bookmark routes
	bookmarkGroup := public.Group("/bookmark", middleware.AuthAny)
	bookmarkGroup.Get("/", handler.GetAllBookmark)
	bookmarkGroup.Post("/", handler.CreateBookmark)
	bookmarkGroup.Get("/acc/:acc_id", handler.GetBookmarkByAccountID)
	bookmarkGroup.Delete("/:id", handler.DeleteBookmark)
	bookmarkGroup.Delete("/ann/:id", handler.DeleteBookmarkByAnnounceID)

	// Answer routes
	answerGroup := public.Group("/answer", middleware.AuthAny)
	answerGroup.Get("/", handler.GetUserQuestions)
	answerGroup.Post("/", handler.CreateUserQuestion)
	answerGroup.Put("/", handler.UpdateUserQuestion)

	// Notification routes
	notificationGroup := public.Group("/notification")
	notificationGroup.Get("/", handler.GetAllNotification)
	notificationGroup.Get("/acc/:acc_id", handler.GetNotificationByAccountID)

	// FCM Token routes
	fcmGroup := public.Group("/fcm")
	fcmGroup.Get("/", handler.GetAllFCMToken)
	fcmGroup.Post("/", handler.CreateFCMToken)
}
