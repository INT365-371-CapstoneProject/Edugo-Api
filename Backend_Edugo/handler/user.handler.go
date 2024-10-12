package handler

import "github.com/gofiber/fiber/v3"

func UserHandlerRead(ctx fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"data": "user"})
}
