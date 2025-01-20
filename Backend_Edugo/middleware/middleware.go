package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func Auth(ctx fiber.Ctx) error {
	token := ctx.Get("x-token")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	// _ , err := utils.VerifyToken(token)
	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	
	role := claims["role"].(string)
	if role != "provider" {
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