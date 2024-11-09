package handler

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/utils"
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

	// Validate Request
	validate := validator.New()
	if errValidate := validate.Struct(post); errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return ctx.Status(400).JSON(fiber.Map{
					"error": "The post type field is required.",
				})
			case "enum":
				return ctx.Status(400).JSON(fiber.Map{
					"error": "The post type must be either 'Announce' or 'Subject'.",
				})
			}
		}
	}

	// Handle File Image
	filename := ctx.Locals("filenameImage").(*string)
	// Handle File Attach
	filenameAttach := ctx.Locals("filenameAttach").(*string)
	// Create New Post
	newPost := entity.Post{
		Title:        post.Title,
		Description:  post.Description,
		URL:          post.URL,
		Attach_File:  filenameAttach,
		Image:        filename,
		Posts_Type:   post.Posts_Type,
		Publish_Date: post.Publish_Date,
		Close_Date:   post.Close_Date,
		Provider_ID:  post.Provider_ID,
		User_ID:      post.User_ID,
	}
	errCreatePost := database.DB.Create(&newPost).Error
	if errCreatePost != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error message": errCreatePost.Error(),
		})
	}

	return ctx.Status(201).JSON(newPost)
}

func GetPostByID(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post entity.Post
	err := database.DB.First(&post, "posts_id = ?", postId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error message": err.Error(),
		})
	}
	return ctx.JSON(post)
}

func UpdatePost(ctx fiber.Ctx) error {
	postRequest := new(request.PostUpdateRequest)
	if err := ctx.Bind().Body(postRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error message": "Bad Request",
		})
	}
	var post entity.Post
	postId := ctx.Params("id")
	// Check Available Post
	err := database.DB.First(&post, "posts_id = ?", postId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error message": "Post not found",
		})
	} else {
		// Update File Image
		// 1.Handle File Image
		_, errFile := ctx.FormFile("image")
		if errFile != nil {
			log.Println("Error File = ", errFile)
		} else {
			// Remove Old File Image
			if post.Image != nil {
				errDeleteImage := utils.HandleRemoveFileImage(*post.Image)
				if errDeleteImage != nil {
					log.Println("Failed to remove image file:", errDeleteImage)
				}
			}
			// Add New File Image
			filename := ctx.Locals("filenameImage").(*string)
			post.Image = filename
		}

		// Update File Attach
		// 2.Handle File Attach
		_, errFileAttach := ctx.FormFile("attach_file")
		if errFileAttach != nil {
			log.Println("Error File Attach = ", errFileAttach)
		} else {
			// Remove Old File Attach
			if post.Attach_File != nil {
				errDeleteAttach := utils.HandleRemoveFileAttach(*post.Attach_File)
				if errDeleteAttach != nil {
					log.Println("Failed to remove attach file:", errDeleteAttach)
				}
			}
			// Add New File Attach
			filenameAttach := ctx.Locals("filenameAttach").(*string)
			post.Attach_File = filenameAttach
		}
		// Update Post
		if postRequest.Title != "" {
			post.Title = postRequest.Title
		}
		if postRequest.Description != "" {
			post.Description = postRequest.Description
		}
		post.URL = postRequest.URL
		post.Close_Date = postRequest.Close_Date

		errUpdate := database.DB.Save(&post).Error
		if errUpdate != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"error message": errUpdate.Error(),
			})
		}
	}
	return ctx.Status(200).JSON(post)
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
	} else {
		// Handle File Remove Image
		if post.Image != nil {
			errDeleteImage := utils.HandleRemoveFileImage(*post.Image)
			if errDeleteImage != nil {
				log.Println("Failed to remove image file:", errDeleteImage)
			}
		}

		// Handle File Remove Attach
		if post.Attach_File != nil {
			errDeleteAttach := utils.HandleRemoveFileAttach(*post.Attach_File)
			if errDeleteAttach != nil {
				log.Println("Failed to remove attach file:", errDeleteAttach)
			}
		}

		errDeletePost := database.DB.Debug().Delete(&post, "posts_id = ?", postId).Error
		if errDeletePost != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"error message": errDeletePost.Error(),
			})
		}
		return ctx.JSON(fiber.Map{
			"message": "post deleted",
		})
	}
}
