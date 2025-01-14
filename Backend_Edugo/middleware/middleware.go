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
	_ , err := utils.VerifyToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	return ctx.Next()
}

func PermissionCreate( ctx fiber.Ctx) error {
	token := ctx.Get("x-token")
	if token != "secret" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	return ctx.Next()
}