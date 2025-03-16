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
	// "gorm.io/gorm"
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
	claims["account_id"] = account.Account_ID // แก้จาก account.ID เป็น account.AccountID
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
	if err := database.DB.First(&account, "account_id = ?", claims["account_id"]).Error; err != nil {
		return utils.HandleError(ctx, 404, "User not found")
	}

	switch account.Role {
	case "superadmin":
		return getAdminProfile(ctx, account)
	case "admin":
		return getAdminProfile(ctx, account)
	case "provider":
		return getProviderProfile(ctx, account)
	case "user":
		return getUserProfile(ctx, account)
	default:
		return utils.HandleError(ctx, 400, "Invalid role")
	}
}

func getAdminProfile(ctx fiber.Ctx, account entity.Account) error {
	// Change this line to use a pointer
	var adminDetails *entity.Admin
	// Use &adminDetails to pass the pointer
	if err := database.DB.First(&adminDetails, "account_id = ?", account.Account_ID).Error; err != nil {
		return utils.HandleError(ctx, 404, "Admin details not found")
	}

	return ctx.JSON(fiber.Map{
		"profile": response.AdminProfileResponse{
			ID:        account.Account_ID,
			Email:     account.Email,
			Username:  account.Username,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Role:      account.Role,
			Phone:     &adminDetails.Phone,
			Status:   account.Status,
		},
	})
}

func getProviderProfile(ctx fiber.Ctx, account entity.Account) error {
	var providerDetails entity.Provider
	if err := database.DB.First(&providerDetails, "account_id = ?", account.Account_ID).Error; err != nil {
		return utils.HandleError(ctx, 404, "Provider details not found")
	}

	// สร้างตัวแปรเพื่อเก็บค่า default สำหรับ fields ที่อาจเป็น nil
	firstName := ""
	if account.FirstName != nil {
		firstName = *account.FirstName
	}

	lastName := ""
	if account.LastName != nil {
		lastName = *account.LastName
	}

	phonePerson := ""
	if providerDetails.Phone_Person != nil {
		phonePerson = *providerDetails.Phone_Person
	}

	return ctx.JSON(fiber.Map{
		"profile": response.ProviderProfileResponse{
			ID:           account.Account_ID,
			Email:        account.Email,
			Username:     account.Username,
			FirstName:    firstName,
			LastName:     lastName,
			Role:         account.Role,
			Company_Name: providerDetails.Company_Name,
			Phone:        providerDetails.Phone,
			Phone_Person: phonePerson,
			Address:      providerDetails.Address,
			City:         providerDetails.City,
			Country:      providerDetails.Country,
			Postal_Code:  providerDetails.Postal_Code,
			Status:       account.Status,
			Verify:       providerDetails.Verify,
		},
	})
}

func getUserProfile(ctx fiber.Ctx, account entity.Account) error {
	var userDetails entity.Account
	if err := database.DB.First(&userDetails, "account_id = ?", account.Account_ID).Error; err != nil {
		return utils.HandleError(ctx, 404, "User details not found")
	}

	return ctx.JSON(fiber.Map{
		"profile": response.UserProfileResponse{
			ID:        account.Account_ID,
			Email:     account.Email,
			Username:  account.Username,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Role:      account.Role,
			Status:    account.Status,
		},
	})
}

func UpdateProfile(ctx fiber.Ctx) error {
	claims := middleware.GetTokenClaims(ctx)
	var account entity.Account
	if err := database.DB.First(&account, "account_id = ?", claims["account_id"]).Error; err != nil {
		return utils.HandleError(ctx, 404, "User not found")
	}

	switch account.Role {
	case "superadmin":
		return updateAdminProfile(ctx, account)
	case "admin":
		return updateAdminProfile(ctx, account)
	case "provider":
		return updateProviderProfile(ctx, account)
	case "user":
		return updateUserProfile(ctx, account)
	default:
		return utils.HandleError(ctx, 400, "Invalid role")
	}
}

func updateAdminProfile(ctx fiber.Ctx, account entity.Account) error {
	// Handle avatar upload first
	if err := utils.HandleAvatarUpload(ctx, "avatar"); err != nil {
		return utils.HandleError(ctx, 400, "Failed to upload avatar")
	}

	avatarBytes := ctx.Locals("avatarBytes")

	updateRequest := new(request.AdminUpdateRequest)
	if err := ctx.Bind().Body(updateRequest); err != nil {
		return utils.HandleError(ctx, 400, "Invalid request body")
	}

	// Validate Request
	if err := validate.Struct(updateRequest); err != nil {
		return utils.HandleError(ctx, 400, "Validation failed")
	}

	// Start transaction
	tx := database.DB.Begin()

	// Build updates map with only provided fields
	updates := make(map[string]interface{})
	if updateRequest.Username != nil {
		updates["username"] = *updateRequest.Username
	}
	if updateRequest.Email != nil {
		updates["email"] = *updateRequest.Email
	}
	if updateRequest.FirstName != nil {
		updates["first_name"] = *updateRequest.FirstName
	}
	if updateRequest.LastName != nil {
		updates["last_name"] = *updateRequest.LastName
	}

	// Add avatar to updates if provided
	if avatarBytes != nil {
		updates["avatar"] = avatarBytes
	}

	// Update account information if there are changes
	if len(updates) > 0 {
		if err := tx.Model(&account).Updates(updates).Error; err != nil {
			tx.Rollback()
			return utils.HandleError(ctx, 500, "Failed to update account")
		}
	}

	// Update admin phone if provided
	if updateRequest.Phone != nil {
		if err := tx.Model(&entity.Admin{}).Where("account_id = ?", account.Account_ID).
			Update("phone", updateRequest.Phone).Error; err != nil {
			tx.Rollback()
			return utils.HandleError(ctx, 500, "Failed to update admin details")
		}
	}

	tx.Commit()

	// Fetch updated details for response
	var adminDetails entity.Admin
	if err := database.DB.Where("account_id = ?", account.Account_ID).First(&adminDetails).Error; err != nil {
		return utils.HandleError(ctx, 404, "Admin details not found")
	}

	return ctx.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"profile": response.AdminProfileResponse{
			ID:        account.Account_ID,
			Email:     account.Email,
			Username:  account.Username,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Role:      account.Role,
			Phone:     &adminDetails.Phone,
		},
	})
}

func updateProviderProfile(ctx fiber.Ctx, account entity.Account) error {
	// Handle avatar upload first
	if err := utils.HandleAvatarUpload(ctx, "avatar"); err != nil {
		return utils.HandleError(ctx, 400, "Failed to upload avatar")
	}

	avatarBytes := ctx.Locals("avatarBytes")

	updateRequest := new(request.ProviderUpdateRequest)
	if err := ctx.Bind().Body(updateRequest); err != nil {
		return utils.HandleError(ctx, 400, "Invalid request body")
	}

	// Validate Request
	if err := validate.Struct(updateRequest); err != nil {
		return utils.HandleError(ctx, 400, "Validation failed")
	}

	// Start transaction
	tx := database.DB.Begin()

	// Build account updates map with only provided fields
	accountUpdates := make(map[string]interface{})
	if updateRequest.Username != nil {
		accountUpdates["username"] = *updateRequest.Username
	}
	if updateRequest.Email != nil {
		accountUpdates["email"] = *updateRequest.Email
	}
	if updateRequest.FirstName != nil {
		accountUpdates["first_name"] = *updateRequest.FirstName
	}
	if updateRequest.LastName != nil {
		accountUpdates["last_name"] = *updateRequest.LastName
	}

	// Add avatar to account updates if provided
	if avatarBytes != nil {
		accountUpdates["avatar"] = avatarBytes
	}

	// Update account information if there are changes
	if len(accountUpdates) > 0 {
		if err := tx.Model(&account).Updates(accountUpdates).Error; err != nil {
			tx.Rollback()
			return utils.HandleError(ctx, 500, "Failed to update account")
		}
	}

	// Build provider updates map with only provided fields
	providerUpdates := make(map[string]interface{})
	if updateRequest.CompanyName != nil {
		providerUpdates["company_name"] = *updateRequest.CompanyName
	}
	if updateRequest.Phone != nil {
		providerUpdates["phone"] = *updateRequest.Phone
	}
	if updateRequest.PhonePerson != nil {
		providerUpdates["phone_person"] = *updateRequest.PhonePerson
	}
	if updateRequest.Address != nil {
		providerUpdates["address"] = *updateRequest.Address
	}
	if updateRequest.City != nil {
		providerUpdates["city"] = *updateRequest.City
	}
	if updateRequest.Country != nil {
		providerUpdates["country"] = *updateRequest.Country
	}
	if updateRequest.PostalCode != nil {
		providerUpdates["postal_code"] = *updateRequest.PostalCode
	}

	// Update provider details if there are changes
	if len(providerUpdates) > 0 {
		if err := tx.Model(&entity.Provider{}).Where("account_id = ?", account.Account_ID).
			Updates(providerUpdates).Error; err != nil {
			tx.Rollback()
			return utils.HandleError(ctx, 500, "Failed to update provider details")
		}
	}

	tx.Commit()

	// Refresh account data after update
	if err := database.DB.First(&account, "account_id = ?", account.Account_ID).Error; err != nil {
		return utils.HandleError(ctx, 404, "Account not found")
	}

	// Fetch updated provider details
	var providerDetails entity.Provider
	if err := database.DB.Where("account_id = ?", account.Account_ID).First(&providerDetails).Error; err != nil {
		return utils.HandleError(ctx, 404, "Provider details not found")
	}

	// Create safe string values for FirstName and LastName
	firstName := ""
	lastName := ""
	if account.FirstName != nil {
		firstName = *account.FirstName
	}
	if account.LastName != nil {
		lastName = *account.LastName
	}

	// สร้างค่า default สำหรับ Phone_Person
	phonePerson := ""
	if providerDetails.Phone_Person != nil {
		phonePerson = *providerDetails.Phone_Person
	}

	// Create response with null checks
	response := response.ProviderProfileResponse{
		ID:           account.Account_ID,
		Email:        account.Email,
		Username:     account.Username,
		FirstName:    firstName, // ใช้ค่าที่ตรวจสอบแล้ว
		LastName:     lastName,  // ใช้ค่าที่ตรวจสอบแล้ว
		Role:         account.Role,
		Company_Name: providerDetails.Company_Name,
		Phone:        providerDetails.Phone,
		Phone_Person: phonePerson, // ใช้ค่าที่ตรวจสอบแล้ว
		Address:      providerDetails.Address,
		City:         providerDetails.City,
		Country:      providerDetails.Country,
		Postal_Code:  providerDetails.Postal_Code,
	}

	return ctx.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"profile": response,
	})
}

func updateUserProfile(ctx fiber.Ctx, account entity.Account) error {
	// Handle avatar upload first
	if err := utils.HandleAvatarUpload(ctx, "avatar"); err != nil {
		return utils.HandleError(ctx, 400, "Failed to upload avatar")
	}

	avatarBytes := ctx.Locals("avatarBytes")

	updateRequest := new(request.UserUpdateRequest)
	if err := ctx.Bind().Body(updateRequest); err != nil {
		return utils.HandleError(ctx, 400, "Invalid request body")
	}

	// Validate Request
	if err := validate.Struct(updateRequest); err != nil {
		return utils.HandleError(ctx, 400, "Validation failed")
	}

	// Build updates map with only provided fields
	updates := make(map[string]interface{})
	if updateRequest.Username != nil {
		updates["username"] = *updateRequest.Username
	}
	if updateRequest.Email != nil {
		updates["email"] = *updateRequest.Email
	}
	if updateRequest.FirstName != nil {
		updates["first_name"] = *updateRequest.FirstName
	}
	if updateRequest.LastName != nil {
		updates["last_name"] = *updateRequest.LastName
	}

	// Add avatar to updates if provided
	if avatarBytes != nil {
		updates["avatar"] = avatarBytes
	}

	// Update account information if there are changes
	if len(updates) > 0 {
		if err := database.DB.Model(&account).Updates(updates).Error; err != nil {
			return utils.HandleError(ctx, 500, "Failed to update profile")
		}
	}

	return ctx.JSON(fiber.Map{
		"profile": response.UserProfileResponse{
			ID:        account.Account_ID,
			Email:     account.Email,
			Username:  account.Username,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Role:      account.Role,
		},
	})
}

func GetAvatarImage(ctx fiber.Ctx) error {
	claims := middleware.GetTokenClaims(ctx)

	var account entity.Account
	if err := database.DB.Select("avatar").First(&account, "account_id = ?", claims["account_id"]).Error; err != nil {
		return utils.HandleError(ctx, 404, "Avatar not found")
	}

	// If no avatar is stored
	if len(account.Avatar) == 0 {
		return utils.HandleError(ctx, 404, "No avatar image found")
	}

	// Set content type header for image
	ctx.Set("Content-Type", "image/jpeg")             // You might want to store the content type in DB if you support multiple formats

	// Return the image bytes directly
	return ctx.Send(account.Avatar)
}
