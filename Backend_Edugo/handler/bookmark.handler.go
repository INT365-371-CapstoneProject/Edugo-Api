package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/middleware"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func CreateBookmark(ctx fiber.Ctx) error {

	claims := middleware.GetTokenClaims(ctx)
	username := claims["username"].(string)

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return handleError(ctx, 404, "Account not found")
	}

	bookmark := new(request.CreateBookmarkRequest)
	if err := ctx.Bind().Body(bookmark); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(bookmark); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// check available post
	var announce entity.Announce_Post
	result := database.DB.Where("announce_id = ?", bookmark.Announce_ID).First(&announce)
	if result.RowsAffected == 0 {
		return utils.HandleError(ctx, 400, "Announce not found")
	}

	// Create bookmark
	newbookmark := entity.Bookmark{
		Announce_ID: bookmark.Announce_ID,
		Account_ID:  account.Account_ID,
	}

	// Create bookmark with debug log
	if err := database.DB.Debug().Create(&newbookmark).Error; err != nil {
		return utils.HandleError(ctx, 400, "Error saving bookmark: "+err.Error())
	}

	// Return response
	bookmarkResponse := response.BookmarkResponse{
		Bookmark_ID: newbookmark.Bookmark_ID,
		Announce_ID: newbookmark.Announce_ID,
		Account_ID:  newbookmark.Account_ID,
		CreatedAt:   newbookmark.CreatedAt,
	}

	return ctx.Status(201).JSON(bookmarkResponse)
}

func GetAllBookmark(ctx fiber.Ctx) error {
	var bookmarks []entity.Bookmark

	// ดึงข้อมูลทั้งหมดจากตาราง bookmark
	if err := database.DB.Find(&bookmarks).Error; err != nil {
		return utils.HandleError(ctx, 500, "Error retrieving bookmarks: "+err.Error())
	}

	// สร้าง response list
	var bookmarkResponse []response.BookmarkResponse
	for _, bookmark := range bookmarks {
		bookmarkResponse = append(bookmarkResponse, response.BookmarkResponse{
			Bookmark_ID: bookmark.Bookmark_ID,
			CreatedAt:   bookmark.CreatedAt,
			Announce_ID: bookmark.Announce_ID,
			Account_ID:  bookmark.Account_ID,
		})
	}

	return ctx.Status(200).JSON(bookmarkResponse)
}

func GetBookmarkByAccountID(ctx fiber.Ctx) error {
	AccountID := ctx.Params("acc_id")
	fmt.Println("AccountID:", AccountID)
	var bookmarks []entity.Bookmark

	// ค้นหาความคิดเห็นที่มี post_id ตรงกับค่าที่ระบุ
	if err := database.DB.Where("account_id = ?", AccountID).Find(&bookmarks).Error; err != nil {
		return utils.HandleError(ctx, 500, "Error retrieving bookmarks: "+err.Error())
	}

	if len(bookmarks) == 0 {
		return utils.HandleError(ctx, 404, "No bookmarks found for this post")
	}

	// สร้าง response list
	var bookmarkResponses []response.BookmarkResponse
	for _, bookmark := range bookmarks {
		bookmarkResponses = append(bookmarkResponses, response.BookmarkResponse{
			Bookmark_ID: bookmark.Bookmark_ID,
			CreatedAt:   bookmark.CreatedAt,
			Announce_ID: bookmark.Announce_ID,
			Account_ID:  bookmark.Account_ID,
		})
	}

	return ctx.Status(200).JSON(bookmarkResponses)
}

func DeleteBookmark(ctx fiber.Ctx) error {
	bookmarkId := ctx.Params("id")

	var bookmark entity.Bookmark
	err := database.DB.Where("bookmark_id = ?", bookmarkId).First(&bookmark).Error
	if err != nil {
		return handleError(ctx, 404, "Bookmark not found")
	}

	err = database.DB.Delete(&bookmark).Error
	if err != nil {
		return handleError(ctx, 400, "Failed to delete bookmark")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Bookmark deleted successfully",
	})
}

func DeleteBookmarkByAnnounceID(ctx fiber.Ctx) error {

	claims := middleware.GetTokenClaims(ctx)
	username := claims["username"].(string)

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return handleError(ctx, 404, "Account not found")
	}

	annId := ctx.Params("id")

	var bookmark entity.Bookmark
	err := database.DB.Where("announce_id = ? And account_id = ?", annId, account.Account_ID).First(&bookmark).Error
	if err != nil {
		return handleError(ctx, 404, "Bookmark not found")
	}

	err = database.DB.Delete(&bookmark).Error
	if err != nil {
		return handleError(ctx, 400, "Failed to delete bookmark")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Bookmark deleted successfully",
	})
}
