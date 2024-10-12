package main

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/route"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var users = []User{
	{ID: 1, FirstName: "John", LastName: "Doe"},
	{ID: 2, FirstName: "Jane", LastName: "Doe"},
	{ID: 3, FirstName: "Joey", LastName: "boy"},
}

func getUsers(c fiber.Ctx) error {
	return c.JSON(users)
}

func getUser(c fiber.Ctx) error {
	id := c.Params("id")
	for _, user := range users {
		if strconv.Itoa(user.ID) == id {
			return c.JSON(user)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
}

func createUser(c fiber.Ctx) error {
	user := new(User)
	if err := c.Bind().Body(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	user.ID = len(users) + 1
	users = append(users, *user)
	return c.JSON(user)
}

func updateUser(c fiber.Ctx) error {
	id := c.Params("id")
	for i, user := range users {
		if strconv.Itoa(user.ID) == id {
			updateUser := new(User)
			if err := c.Bind().Body(updateUser); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}
			updateUser.ID = user.ID
			users[i] = *updateUser
			return c.JSON(updateUser)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
}

func deleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	for i, user := range users {
		if strconv.Itoa(user.ID) == id {
			users = append(users[:i], users[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
}

func main() {
	app := fiber.New()

	// Initial route
	route.RouteInit(app)

	app.Listen(":8080")
}
