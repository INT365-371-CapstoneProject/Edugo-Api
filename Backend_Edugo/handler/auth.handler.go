package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func Login(ctx fiber.Ctx) error {

	loginRequest := new(request.LoginRequest)
	if err := ctx.Bind().Body(loginRequest); err != nil {
		return utils.HandleError(ctx, 400, "Invalid request body")
	}

	// Validate Request
	if errValidate := validate.Struct(loginRequest); errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			return utils.HandleError(ctx, 400, err.Translate(trans))
		}
	}

	// Check Available User
	var account entity.Account
	err := database.DB.First(&account, "email = ? OR username = ?", loginRequest.Email, loginRequest.Username).Error
	if err != nil {
		return utils.HandleError(ctx, 404, "User not found")
	}

	// Check Validated Password
	isValid := utils.CheckPasswordHash(loginRequest.Password, account.Password)
	if !isValid {
		return utils.HandleError(ctx, 400, "Invalid password")
	}

	// กำหนดระยะเวลา token ตาม RememberMe
	var expirationTime time.Time
	if loginRequest.RememberMe {
		expirationTime = time.Now().Add(time.Hour * 24 * 30) // 30 days
	} else {
		expirationTime = time.Now().Add(time.Hour * 24) // 1 day
	}

	// Generate JWT Token
	claims := jwt.MapClaims{}
	claims["email"] = account.Email
	claims["username"] = account.Username
	claims["role"] = account.Role
	// Set token expire time to 7 days
	claims["exp"] = expirationTime.Unix()

	token, errGenerateToken := utils.GenerateToken(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		return utils.HandleError(ctx, 500, "Failed to generate token")
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func ForgotPassword(ctx fiber.Ctx) error {
	forgotRequest := new(request.ForgotPasswordRequest)
	if err := ctx.Bind().Body(forgotRequest); err != nil {
		return utils.HandleError(ctx, 400, "Invalid request body")
	}

	// Validate Request
	if errValidate := validate.Struct(forgotRequest); errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			return utils.HandleError(ctx, 400, err.Translate(trans))
		}
	}

	// Check if user exists
	var account entity.Account
	if err := database.DB.First(&account, "email = ?", forgotRequest.Email).Error; err != nil {
		return utils.HandleError(ctx, 404, "User not found")
	}

	// Generate OTP
	otpCode := utils.GenerateOTP()
	expiredAt := time.Now().Add(15 * time.Minute)

	// Save OTP to database
	otp := entity.OTP{
		Code:      otpCode,
		ExpiredAt: expiredAt,
		AccountID: account.Account_ID, // แก้จาก account.ID เป็น account.AccountID
	}

	if err := database.DB.Create(&otp).Error; err != nil {
		return utils.HandleError(ctx, 500, "Failed to generate OTP")
	}

	// Send OTP via email
	if err := utils.SendOTPEmail(account.Email, otpCode); err != nil {
		log.Printf("Failed to send OTP email: %v", err)
		return utils.HandleError(ctx, 500, fmt.Sprintf("Failed to send OTP email: %v", err))
	}

	return ctx.JSON(fiber.Map{
		"message": "OTP has been sent to your email",
	})
}

func VerifyOTP(ctx fiber.Ctx) error {
	verifyRequest := new(request.VerifyOTPRequest)
	if err := ctx.Bind().Body(verifyRequest); err != nil {
		return utils.HandleError(ctx, 400, "Invalid request body")
	}

	// Validate Request
	if errValidate := validate.Struct(verifyRequest); errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			return utils.HandleError(ctx, 400, err.Translate(trans))
		}
	}

	// Find account
	var account entity.Account
	if err := database.DB.First(&account, "email = ?", verifyRequest.Email).Error; err != nil {
		return utils.HandleError(ctx, 404, "User not found")
	}

	// Find latest valid OTP
	var otp entity.OTP
	if err := database.DB.Where("account_id = ? AND is_used = ? AND attempt_count < ? AND expired_at > ?",
		account.Account_ID, false, 3, time.Now()). // แก้จาก account.ID เป็น account.AccountID
		Order("created_at DESC").
		First(&otp).Error; err != nil {
		return utils.HandleError(ctx, 400, "No valid OTP found")
	}

	// Validate OTP
	if !utils.ValidateOTP(otp.Code, verifyRequest.OTPCode, otp.ExpiredAt) {
		// Increment attempt count
		otp.AttemptCount++
		database.DB.Save(&otp)
		return utils.HandleError(ctx, 400, "Invalid OTP")
	}

	// Hash new password
	hashedPassword, err := utils.HashingPassword(verifyRequest.NewPassword)
	if err != nil {
		return utils.HandleError(ctx, 500, "Failed to hash password")
	}

	// Update password and mark OTP as used
	tx := database.DB.Begin()
	
	if err := tx.Model(&account).Update("password", hashedPassword).Error; err != nil {
		tx.Rollback()
		return utils.HandleError(ctx, 500, "Failed to update password")
	}

	if err := tx.Model(&otp).Update("is_used", true).Error; err != nil {
		tx.Rollback()
		return utils.HandleError(ctx, 500, "Failed to update OTP status")
	}

	tx.Commit()

	return ctx.JSON(fiber.Map{
		"message": "Password has been reset successfully",
	})
}
