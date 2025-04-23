package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/middleware"
	"github.com/tk-neng/demo-go-fiber/model/entity"

	// "github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func GetAllNotification(ctx fiber.Ctx) error {
	var notifications []entity.Notification

	// ดึงข้อมูลทั้งหมดจากตาราง bookmark
	if err := database.DB.Find(&notifications).Error; err != nil {
		return utils.HandleError(ctx, 500, "Error retrieving notifications: "+err.Error())
	}

	// สร้าง response list
	var notificationResponse []response.NotificationResponse
	for _, notification := range notifications {
		notificationResponse = append(notificationResponse, response.NotificationResponse{
			Notification_ID: notification.Notification_ID,
			Title:           notification.Title,
			Message:         notification.Message,
			Created_At:      notification.Created_At,
			Is_Read:         notification.Is_Read,
			Announce_ID:     notification.Announce_ID,
			Account_ID:      notification.Account_ID,
		})
	}

	return ctx.Status(200).JSON(notificationResponse)
}

func GetNotificationByAccountID(ctx fiber.Ctx) error {
	AccountID := ctx.Params("acc_id")
	fmt.Println("AccountID:", AccountID)
	var notifications []entity.Notification

	// ค้นหาความคิดเห็นที่มี notification_id ตรงกับค่าที่ระบุ
	if err := database.DB.Where("account_id = ?", AccountID).Find(&notifications).Error; err != nil {
		return utils.HandleError(ctx, 500, "Error retrieving notifications: "+err.Error())
	}

	if len(notifications) == 0 {
		return utils.HandleError(ctx, 404, "No notifications found")
	}

	// สร้าง response list
	var notificationResponse []response.NotificationResponse
	for _, notification := range notifications {
		notificationResponse = append(notificationResponse, response.NotificationResponse{
			Notification_ID: notification.Notification_ID,
			Title:           notification.Title,
			Message:         notification.Message,
			Created_At:      notification.Created_At,
			Is_Read:         notification.Is_Read,
			Announce_ID:     notification.Announce_ID,
			Account_ID:      notification.Account_ID,
		})
	}

	return ctx.Status(200).JSON(notificationResponse)
}

func UpdateNotification(ctx fiber.Ctx) error {
	claims := middleware.GetTokenClaims(ctx)
	username := claims["username"].(string)

	notificationId := ctx.Params("id")

	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	var notification entity.Notification
	err := database.DB.Where("notification_id = ? AND account_id = ?",
		notificationId, account.Account_ID).First(&notification).Error
	if err != nil {
		fmt.Println("Account ID:", account.Account_ID)
		return handleError(ctx, 403, "Forbidden")
	}

	// Begin a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return handleError(ctx, 409, "Failed to begin transaction")
	}

	notification.Is_Read = 0

	// Save updated notification record
	if err := tx.Save(&notification).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "Failed to update notification details")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "Failed to commit transaction")
	}

	// Construct response data
	notificationResponse := response.NotificationUpdateResponse{
		Notification_ID: notification.Notification_ID,
		Is_Read:         notification.Is_Read,
		Announce_ID:     notification.Announce_ID,
		Account_ID:      account.Account_ID,
	}

	// Return the updated response
	return ctx.Status(200).JSON(notificationResponse)
}
