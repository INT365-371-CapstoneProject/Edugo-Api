package handler

import (

	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/response"
)

func GetAllUser(ctx fiber.Ctx) error {
	var users []entity.User
	result := database.DB.Preload("Account", "role = ?", "user").Find(&users)
	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}
	
	var userResponse []response.UserResponse
	for _, user := range users {
		userResponse = append(userResponse, response.UserResponse{
			User_ID:     user.User_ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Phone_Number: user.Account.Phone_Number,
			Update_On:   user.Account.Update_On,
			Last_Login:  user.Account.Last_Login,
			Username:    user.Account.Username,
			Email:       user.Account.Email,
			Role:        user.Account.Role,
		})
	}
	return ctx.Status(200).JSON(userResponse)
}