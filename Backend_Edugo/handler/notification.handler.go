package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
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
			Notification_ID: notification.NotificationID,
			Title:           notification.Title,
			Message:         notification.Message,
			CreatedAt:       notification.CreatedAt,
			IsRead:          notification.IsRead,
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

	// ค้นหาความคิดเห็นที่มี post_id ตรงกับค่าที่ระบุ
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
			Notification_ID: notification.NotificationID,
			Title:           notification.Title,
			Message:         notification.Message,
			CreatedAt:       notification.CreatedAt,
			IsRead:          notification.IsRead,
			Announce_ID:     notification.Announce_ID,
			Account_ID:      notification.Account_ID,
		})
	}

	return ctx.Status(200).JSON(notificationResponse)
}
