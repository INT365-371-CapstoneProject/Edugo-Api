package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
)

var im int = 1
var pdf int = 1

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
	file, errFile := ctx.FormFile("image")
	if errFile != nil {
		log.Println("Error File = ", errFile)
	}

	var filename *string
	if file != nil {
		filename = &file.Filename
		extenstionFile := filepath.Ext(*filename)
		newFilename := fmt.Sprintf("%d%s", im, extenstionFile)
		errSaveFile := ctx.SaveFile(file, fmt.Sprintf("./public/images/%s", newFilename))
		if errSaveFile != nil {
			log.Println("Fail to store file into public/images directory.")
		} else {
			im++
			filename = &newFilename
		}
	} else {
		log.Println("No file uploaded")
		filename = nil
	}

	// Handle File Attach
	fileAttach, errFileAttach := ctx.FormFile("attach_file")
	if errFileAttach != nil {
		log.Println("Error File Attach = ", errFileAttach)
	}

	var filenameAttach *string
	if fileAttach != nil {
		filenameAttach = &fileAttach.Filename
		extenstionFileAttach := filepath.Ext(*filenameAttach)
		newFilenameAttach := fmt.Sprintf("%d%s", pdf, extenstionFileAttach)
		errSaveFileAttach := ctx.SaveFile(fileAttach, fmt.Sprintf("./public/pdfs/%s", newFilenameAttach))
		if errSaveFileAttach != nil {
			log.Println("Fail to store file into public/attach directory.")
		} else {
			pdf++
			filenameAttach = &newFilenameAttach
		}
	} else {
		log.Println("No file uploaded")
		filenameAttach = nil
	}

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
	}

	// Update File Image
	// 1.Handle File Image
	file, errFile := ctx.FormFile("image")
	if errFile != nil {
		log.Println("Error File = ", errFile)
	} else {
		// Remove Old File Image
		errDeleteImage := os.Remove(fmt.Sprintf("./public/images/%s", *post.Image))
		if errDeleteImage != nil {
			log.Println("Failed to remove image file")
		}
		// Add New File Image
		var filename *string
		if file != nil {
			filename = &file.Filename
			extenstionFile := filepath.Ext(*filename)
			newFilename := fmt.Sprintf("%d%s", im, extenstionFile)
			errSaveFile := ctx.SaveFile(file, fmt.Sprintf("./public/images/%s", newFilename))
			if errSaveFile != nil {
				log.Println("Fail to store file into public/images directory.")
			} else {
				im++
				filename = &newFilename
			}
		} else {
			log.Println("No file uploaded")
			filename = nil
		}
		post.Image = filename
	}

	// Update File Attach
	// 2.Handle File Attach
	fileAttach, errFileAttach := ctx.FormFile("attach_file")
	if errFileAttach != nil {
		log.Println("Error File Attach = ", errFileAttach)
	} else {
		// Remove Old File Attach
		errDeleteAttach := os.Remove(fmt.Sprintf("./public/pdfs/%s", *post.Attach_File))
		if errDeleteAttach != nil {
			log.Println("Failed to remove attach file")
		}
		// Add New File Attach
		var filenameAttach *string
		if fileAttach != nil {
			filenameAttach = &fileAttach.Filename
			extenstionFileAttach := filepath.Ext(*filenameAttach)
			newFilenameAttach := fmt.Sprintf("%d%s", pdf, extenstionFileAttach)
			errSaveFileAttach := ctx.SaveFile(fileAttach, fmt.Sprintf("./public/pdfs/%s", newFilenameAttach))
			if errSaveFileAttach != nil {
				log.Println("Fail to store file into public/attach directory.")
			} else {
				pdf++
				filenameAttach = &newFilenameAttach
			}
		} else {
			log.Println("No file uploaded")
			filenameAttach = nil
		}
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

	// Remove File Image
	errDeleteImage := os.Remove(fmt.Sprintf("./public/images/%s", *post.Image))
	if errDeleteImage != nil {
		log.Println("Failed to remove image file")
	}

	// Remove File Attach
	errDeleteAttach := os.Remove(fmt.Sprintf("./public/pdfs/%s", *post.Attach_File))
	if errDeleteAttach != nil {
		log.Println("Failed to remove attach file")
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
