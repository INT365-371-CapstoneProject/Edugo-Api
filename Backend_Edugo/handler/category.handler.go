package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
)

func GetAllCategory(ctx fiber.Ctx) error {
	var categories []entity.Category
	result := database.DB.Order("Category_ID asc").Find(&categories)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error message": result.Error.Error(),
		})
	}
	return ctx.Status(200).JSON(categories)
}