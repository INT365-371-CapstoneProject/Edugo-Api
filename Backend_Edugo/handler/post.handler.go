package handler

import (
	"time"

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
	if errValidate := validate.Struct(post); errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return ctx.Status(400).JSON(fiber.Map{
					"error": "The post type field is required.",
				})
			case "oneof":
				return ctx.Status(400).JSON(fiber.Map{
					"error": "The post type must be either 'Announce' or 'Subject'.",
				})
			}
		}
	}
	if post.Publish_Date == nil {
		now := time.Now().UTC()
		post.Publish_Date = &now
	}else{
		utcTime := post.Publish_Date.UTC()
		post.Publish_Date = &utcTime
	}
	if post.Close_Date != nil {
		utcTime := post.Close_Date.UTC()
		post.Close_Date = &utcTime
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

func DeletePost(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post entity.Post
	// Check Available Post
	err := database.DB.Debug().First(&post, "posts_id = ?", postId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "post not found",
		})
	}

	errDeletePost := database.DB.Debug().Delete(&post, "posts_id = ?", postId ).Error
	if errDeletePost != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error message": errDeletePost.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"message": "post deleted",
	})
}