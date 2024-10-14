package handler

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
)

func GetAllPost(ctx fiber.Ctx) error {
	var posts []entity.Post
	result := database.DB.Find(&posts)
	if result.Error != nil {
		log.Println(result.Error)
	}
	return ctx.JSON(posts)
}