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
		Status:     "Active",
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
		Phone:     newAdmin.Phone,
	}

	return ctx.Status(201).JSON(providerResponse)
}
