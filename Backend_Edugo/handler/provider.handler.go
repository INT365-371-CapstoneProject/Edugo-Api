package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
	"gorm.io/gorm"
)

func GetIDProvider(ctx fiber.Ctx) error {
	// Get provider ID from params
	providerID := ctx.Params("id")

	// Find provider in database with Account preloaded
	var provider entity.Provider
	result := database.DB.Preload("Account", "role = ?", "provider").First(&provider, providerID)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Provider not found",
		})
	}

	// Create response
	providerResponse := response.ProviderResponse{
		Provider_ID:  provider.Provider_ID,
		Company_Name: provider.Company_Name,
		FirstName:    provider.Account.FirstName,
		LastName:     provider.Account.LastName,
		Username:     provider.Account.Username,
		Email:        provider.Account.Email,
		URL:          provider.URL,
		Address:      provider.Address,
		City:         provider.City,
		Country:      provider.Country,
		Postal_Code:  provider.Postal_Code,
		Phone:        provider.Phone,
		Status:       provider.Account.Status, // Correctly fetch Status from the preloaded Account
		Verify:       provider.Verify,
		Create_On:    provider.Account.Create_On,
		Last_Login:   provider.Account.Last_Login,
		Update_On:    provider.Account.Update_On,
		Role:         provider.Account.Role,
	}

	return ctx.Status(200).JSON(providerResponse)
}

func CreateProvider(ctx fiber.Ctx) error {
	provider := new(request.ProviderCreateRequest)
	if err := ctx.Bind().Body(provider); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(provider); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// check duplicate email
	var account entity.Account
	result := database.DB.Where("email = ?", provider.Email).First(&account)
	if result.RowsAffected > 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// check duplicate username
	result = database.DB.Where("username = ?", provider.Username).First(&account)
	if result.RowsAffected > 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Username already exists",
		})
	}

	// Begin transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to begin transaction",
		})
	}

	// Create account
	newAccount := entity.Account{
		Username:   provider.Username,
		Email:      provider.Email,
		FirstName:  &provider.FirstName,
		LastName:   &provider.LastName,
		Status:     "Active",
		Last_Login: nil,
		Role:       "provider",
	}

	// Hash password
	hashedPassword, err := utils.HashingPassword(provider.Password)
	if err != nil {
		tx.Rollback()
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}
	newAccount.Password = hashedPassword

	// Insert account to database
	if err := tx.Create(&newAccount).Error; err != nil {
		tx.Rollback()
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to create account",
		})
	}

	// Create provider
	newProvider := entity.Provider{
		Company_Name: provider.Company_name,
		URL:          provider.URL,
		Address:      provider.Address,
		City:         provider.City,
		Country:      provider.Country,
		Postal_Code:  provider.Postal_code,
		Phone:        provider.Phone,
		Verify:       "Waiting",
		Account_ID:   newAccount.Account_ID,
	}

	// Insert provider to database
	if err := tx.Create(&newProvider).Error; err != nil {
		tx.Rollback()
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to create provider",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to commit transaction",
		})
	}

	// Create response
	providerResponse := response.ProviderResponse{
		Provider_ID:  newProvider.Provider_ID,
		Company_Name: newProvider.Company_Name,
		FirstName:    newAccount.FirstName,
		LastName:     newAccount.LastName,
		Username:     newAccount.Username,
		Email:        newAccount.Email,
		URL:          newProvider.URL,
		Address:      newProvider.Address,
		City:         newProvider.City,
		Country:      newProvider.Country,
		Postal_Code:  newProvider.Postal_Code,
		Phone:        newProvider.Phone,
		Status:       newAccount.Status,
		Verify:       newProvider.Verify,
		Create_On:    newAccount.Create_On,
		Last_Login:   newAccount.Last_Login,
		Update_On:    newAccount.Update_On,
		Role:         newAccount.Role,
	}

	return ctx.Status(201).JSON(providerResponse)
}

// GetAvatarImageByProviderID - ดึงรูปภาพ Avatar ของ Provider ตาม ID
func GetAvatarImageByProviderID(ctx fiber.Ctx) error {
	providerID := ctx.Params("id")

	// Find provider to get Account_ID
	var provider entity.Provider
	if err := database.DB.Select("account_id").First(&provider, "provider_id = ?", providerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(ctx, 404, "Provider not found")
		}
		return utils.HandleError(ctx, 500, "Database error finding provider")
	}

	// Find account to get Avatar using Account_ID
	var account entity.Account
	if err := database.DB.Select("avatar").First(&account, "account_id = ?", provider.Account_ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// This case might indicate inconsistent data if the provider exists but the account doesn't
			return utils.HandleError(ctx, 404, "Associated account not found")
		}
		return utils.HandleError(ctx, 500, "Database error finding account")
	}

	// If no avatar is stored
	if len(account.Avatar) == 0 {
		return utils.HandleError(ctx, 404, "No avatar image found for this provider")
	}

	// Set content type header for image
	// Consider storing the actual content type if supporting multiple formats
	ctx.Set("Content-Type", "image/jpeg")

	// Return the image bytes directly
	return ctx.Send(account.Avatar)
}
