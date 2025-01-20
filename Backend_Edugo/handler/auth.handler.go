package handler

import (
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
