package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/utils"
	"strings" // added import
)

func AuthProvider(ctx fiber.Ctx) error {
	token := ctx.Get("Authorization")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	token = parts[1]
	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	
	role := claims["role"].(string)
	if role != "provider" && role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden access",
		})
	}

	// ctx.Locals("user", claims)
	// ctx.Locals("role", claims["role"])
	return ctx.Next()
}

func PermissionCreate( ctx fiber.Ctx) error {
	return ctx.Next()
}