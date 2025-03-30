package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/middleware"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func GetAllFCMToken(ctx fiber.Ctx) error {
	var fcmTokens []entity.FCMToken

	// ดึงข้อมูลทั้งหมดจากตาราง fcmToken
	if err := database.DB.Find(&fcmTokens).Error; err != nil {
		return utils.HandleError(ctx, 500, "Error retrieving fcmTokens: "+err.Error())
	}

	// สร้าง response list
	var fCMTokenResponses []response.FCMTokenResponse
	for _, fCMToken := range fcmTokens {
		fCMTokenResponses = append(fCMTokenResponses, response.FCMTokenResponse{
			Token_ID:   fCMToken.Token_ID,
			FCM_Token:  fCMToken.FCM_Token,
			Created_At: fCMToken.Created_At,
			Account_ID: fCMToken.Account_ID,
		})
	}

	return ctx.Status(200).JSON(fCMTokenResponses)
}

func CreateFCMToken(ctx fiber.Ctx) error {

	claims := middleware.GetTokenClaims(ctx)
	username := claims["username"].(string)

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return handleError(ctx, 404, "Account not found")
	}

	fcmToken := new(request.CreateFCMTokenRequest)
	if err := ctx.Bind().Body(fcmToken); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(fcmToken); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// Create fcmToken
	newFCMToken := entity.FCMToken{
		FCM_Token:  fcmToken.FCM_Token,
		Account_ID: account.Account_ID,
	}

	// Create fcmToken with debug log
	if err := database.DB.Debug().Create(&newFCMToken).Error; err != nil {
		return utils.HandleError(ctx, 400, "Error saving fcmToken: "+err.Error())
	}

	// Return response
	fcmTokenResponse := response.FCMTokenResponse{
		Token_ID:   newFCMToken.Token_ID,
		FCM_Token:  newFCMToken.FCM_Token,
		Account_ID: newFCMToken.Account_ID,
		Created_At: newFCMToken.Created_At,
	}

	return ctx.Status(201).JSON(fcmTokenResponse)
}
