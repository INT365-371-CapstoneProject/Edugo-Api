package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func CreateAdminForSuperadmin(ctx fiber.Ctx) error {
	admin := new(request.AdminCreateRequest)
	if err := ctx.Bind().Body(admin); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(admin); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// check duplicate email
	var account entity.Account
	result := database.DB.Where("email = ?", admin.Email).First(&account)
	if result.RowsAffected > 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// check duplicate username
	result = database.DB.Where("username = ?", admin.Username).First(&account)
	if result.RowsAffected > 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Username already exists",
		})
	}

	// Begin transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return ctx.Status(409).JSON(fiber.Map{
			"message": "Failed to begin transaction",
		})
	}

	// Create account
	newAccount := entity.Account{
		Username:   admin.Username,
		Email:      admin.Email,
		FirstName:  &admin.FirstName,
		LastName:   &admin.LastName,
		Status:     "Active",
		Last_Login: nil,
		Role:       "admin",
	}

	// Hash password
	hashedPassword, err := utils.HashingPassword(admin.Password)
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
		return ctx.Status(409).JSON(fiber.Map{
			"message": "Failed to create account",
		})
	}

	// Create admin - แก้ไขส่วนนี้
	newAdmin := entity.Admin{
		Account_ID: newAccount.Account_ID, // ใช้ Account_ID จาก account ที่เพิ่งสร้าง
		Phone:      admin.Phone,
	}

	// Insert admin to database
	if err := tx.Create(&newAdmin).Error; err != nil {
		tx.Rollback()
		return ctx.Status(409).JSON(fiber.Map{
			"message": "Failed to create admin",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return ctx.Status(409).JSON(fiber.Map{
			"message": "Failed to commit transaction",
		})
	}

	// Create response
	providerResponse := response.AdminResponse{
		ID:        newAccount.Account_ID,
		Username:  newAccount.Username,
		Email:     newAccount.Email,
		FirstName: *newAccount.FirstName,
		LastName:  *newAccount.LastName,
		Status:    newAccount.Status,
		Phone:     newAdmin.Phone,
	}

	return ctx.Status(201).JSON(providerResponse)
}

func GetAllProviderForAdmin(ctx fiber.Ctx) error {
	var providers []entity.Provider
	result := database.DB.Preload("Account", "role = ?", "provider").Find(&providers)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	var providerResponse []response.ProviderResponse
	for _, provider := range providers {
		providerResponse = append(providerResponse, response.ProviderResponse{
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
			Status:       provider.Account.Status,
			Verify:       provider.Verify,
			Create_On:    provider.Account.Create_On,
			Last_Login:   provider.Account.Last_Login,
			Update_On:    provider.Account.Update_On,
			Role:         provider.Account.Role,
		})
	}
	return ctx.Status(200).JSON(providerResponse)
}

func GetIDProviderForAdmin(ctx fiber.Ctx) error {
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
		Status:       provider.Account.Status,
		Verify:       provider.Verify,
		Create_On:    provider.Account.Create_On,
		Last_Login:   provider.Account.Last_Login,
		Update_On:    provider.Account.Update_On,
		Role:         provider.Account.Role,
	}

	return ctx.Status(200).JSON(providerResponse)
}

func VerifyProviderForAdmin(ctx fiber.Ctx) error {
	// Get provider ID from params
	providerID := ctx.Params("id")

	// Bind request body to get verification status
	verifyRequest := struct {
		Status string `json:"status" validate:"required,oneof=Yes No"`
	}{}

	if err := ctx.Bind().Body(&verifyRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request body: " + err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(verifyRequest); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// Find provider in database
	var provider entity.Provider
	result := database.DB.First(&provider, providerID)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Provider not found",
		})
	}

	// ใช้ค่าที่รับมาจาก frontend โดยตรง เนื่องจากค่า ENUM ในฐานข้อมูลตรงกับค่าที่รับมา
	// Update provider verify status
	if err := database.DB.Model(&provider).Update("verify", verifyRequest.Status).Error; err != nil {
		return ctx.Status(409).JSON(fiber.Map{
			"message": "Failed to update provider verification status",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Provider verification status updated to " + verifyRequest.Status,
	})
}

