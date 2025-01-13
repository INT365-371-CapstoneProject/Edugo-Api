package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
)

func Login(ctx fiber.Ctx) error {
	// เพิ่ม log ดู raw body
	fmt.Printf("Raw Request Body: %s\n", string(ctx.Body()))
	
	loginRequest := new(request.LoginRequest)
	if err := ctx.Bind().JSON(loginRequest); err != nil {
		return handleError(ctx, 400, "Invalid request body")
	}

	// Validate Request
	if errValidate := validate.Struct(loginRequest); errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			return handleError(ctx, 400, err.Translate(trans))
		}
	}

	var account entity.Account
	err := database.DB.First(&account, "email = ? OR username = ?", loginRequest.Email, loginRequest.Username).Error
	if err != nil {
		return handleError(ctx, 404, "User not found")
	}

	return ctx.JSON(fiber.Map{
		"token": "secret",
	})
}
