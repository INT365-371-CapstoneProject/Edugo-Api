package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
)

func GatAllCountry(ctx fiber.Ctx) error {
	var countries []entity.Country
	result := database.DB.Order("Country_ID asc").Find(&countries)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error message": result.Error.Error(),
		})
	}
	return ctx.Status(200).JSON(countries)
}