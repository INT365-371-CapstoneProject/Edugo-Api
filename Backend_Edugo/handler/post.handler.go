package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
)

func GetAllPost(ctx fiber.Ctx) error {
	var posts []entity.Post
	result := database.DB.Find(&posts)
	if result.Error != nil {
		// return status 404
		return ctx.Status(404).JSON(fiber.Map{
			"error message": result.Error.Error(),
		})
	}
	return ctx.JSON(posts)
}

func CreatePost(ctx fiber.Ctx) error {
	post := new(request.PostCreateRequest)
	if err := ctx.Bind().Body(post); err != nil {
		return err
	}
	validate := validator.New()
	errValidate := validate.Struct(post)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error message": errValidate.Error(),
		})
	}
	newPost := entity.Post{
		Title: post.Title,
		Description: post.Description,
		URL: post.URL,
		Attach_File: post.Attach_File,
		Posts_Type: post.Posts_Type,
		Publish_Date: post.Publish_Date,
		Close_Date: post.Close_Date,
		Provider_ID: post.Provider_ID,
		User_ID: post.User_ID,
	}
	errCreatePost := database.DB.Create(&newPost).Error
	if errCreatePost != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error message": errCreatePost.Error(),
		})
	}

	return ctx.JSON(newPost)
}

func GetPostByID(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post entity.Post
	err := database.DB.First(&post,"posts_id = ?", postId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error message": err.Error(),
		})
	}
	return ctx.JSON(post)
}