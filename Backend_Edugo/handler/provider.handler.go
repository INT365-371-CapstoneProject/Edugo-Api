package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/response"
)

func GetAllProvider(ctx fiber.Ctx) error {
	var providers []entity.Provider
	result := database.DB.Preload("Account", "role = ?", "provider").Find(&providers)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	var providerResponse []response.ProviderResponse
	for _, provider := range providers {
		providerResponse = append(providerResponse, response.ProviderResponse{
			Provider_ID: provider.Provider_ID,
			Company_Name: provider.Company_Name,
			Username: provider.Account.Username,
			Email: provider.Account.Email,
			URL: provider.URL,
			Address: provider.Address,
			Phone: provider.Phone,
			Status: provider.Status,
			Verify: provider.Verify,
			Create_On: provider.Account.Create_On,
			Last_Login: provider.Account.Last_Login,
			Update_On: provider.Account.Update_On,
			Role: provider.Account.Role,

		})
	}
	return ctx.Status(200).JSON(providerResponse)
}