package handler

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
)

func GetAllUser(ctx fiber.Ctx) error {
	var users []entity.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		log.Println(result.Error)
	}
	return ctx.JSON(users)
}

func CreateUser(ctx fiber.Ctx) error {
	user := new(request.UserCreateRequest)

	if err := ctx.Bind().Body(user); err != nil {
		return err
	}
	newUser := entity.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		AccountID: user.AccountID,
	}

	errCreateUser := database.DB.Create(&newUser).Error

	if errCreateUser != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "User created successfully",
		"data":    newUser,
	})
}


