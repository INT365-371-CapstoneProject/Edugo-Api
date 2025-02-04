package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/middleware"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
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

func GetProfile(ctx fiber.Ctx) error {
	claims := middleware.GetTokenClaims(ctx)
	var account entity.Account
	if err := database.DB.First(&account, "email = ?", claims["email"]).Error; err != nil {
		return utils.HandleError(ctx, 404, "User not found")
	}

	// Response
	profile := response.ProfileResponse{
		ID:       account.Account_ID, // แก้จาก account.ID เป็น account.Account_ID
		Email:    account.Email,
		Username: account.Username,
		Role:     account.Role,
	}
	return ctx.JSON(profile)
}

func EditProfile(ctx fiber.Ctx) error {
    // 1. รับและตรวจสอบ request
    editRequest := new(request.EditProfileRequest)
    if err := ctx.Bind().Body(editRequest); err != nil {
        return utils.HandleError(ctx, 400, "Invalid request body")
    }

    // 2. ดึงข้อมูล account จาก token
    claims := middleware.GetTokenClaims(ctx)
    var account entity.Account
    if err := database.DB.First(&account, "email = ?", claims["email"]).Error; err != nil {
        return utils.HandleError(ctx, 404, "User not found")
    }

    // 3. สร้าง map สำหรับเก็บข้อมูลที่จะอัพเดท
    updates := buildUpdatesMap(editRequest)
    
    // 4. อัพเดทข้อมูลถ้ามีการเปลี่ยนแปลง
    if len(updates) > 0 {
        if err := updateProfile(&account, updates); err != nil {
            return utils.HandleError(ctx, 500, "Failed to update profile")
        }
    }

    // 5. ส่งข้อมูลกลับ
    return ctx.JSON(fiber.Map{
        "message": "Profile updated successfully",
        "profile": response.ProfileResponse{
            ID:        account.Account_ID,
            Email:     account.Email,
            Username:  account.Username,
            FirstName: account.FirstName,
            LastName:  account.LastName,
            Role:      account.Role,
        },
    })
}

// helper functions
func buildUpdatesMap(req *request.EditProfileRequest) map[string]interface{} {
    updates := make(map[string]interface{})
    
    if req.FirstName != nil {
        updates["first_name"] = req.FirstName
    }
    if req.LastName != nil {
        updates["last_name"] = req.LastName
    }
    if req.Email != nil {
        updates["email"] = req.Email
    }
    if req.Username != nil {
        updates["username"] = req.Username
    }
    
    return updates
}

func updateProfile(account *entity.Account, updates map[string]interface{}) error {
    result := database.DB.Model(account).Updates(updates)
    if result.Error != nil {
        log.Printf("Error updating profile: %v", result.Error)
        return result.Error
    }
    
    // ดึงข้อมูลล่าสุด
    return database.DB.First(account, account.Account_ID).Error
}
