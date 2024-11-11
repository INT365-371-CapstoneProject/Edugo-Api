package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func GetAllAnnoucePost(ctx fiber.Ctx) error {
	var posts []entity.Announce_Post
	result := database.DB.Find(&posts)
	if result.Error != nil {
		// return status 404
		return ctx.Status(404).JSON(fiber.Map{
			"error message": result.Error.Error(),
		})
	} else {
		// ใช้ฟังก์ชัน GetCategoryName จากไฟล์ utils/category.go
		utils.GetCategoryName(posts)
		utils.GetPostByAnnounceID(posts)
		// สร้างตัวแปรแบบ slice ของ AnnouncePostResponse
		var postsResponse []response.AnnouncePostResponse
		// วนลูปข้อมูลใน posts
		for _, post := range posts {
			// สร้างตัวแปรแบบ AnnouncePostResponse
			// และกำหนดค่าให้กับตัวแปรนี้
			postsResponse = append(postsResponse, response.AnnouncePostResponse{
				Announce_ID:    post.Announce_ID,
				Title:          post.Post.Title,
				Description:    post.Post.Description,
				URL:            post.Url,
				Attach_File:    post.Attach_File,
				Image:          post.Post.Image,
				Post_Type:      post.Post.Posts_Type,
				Published_Date: post.Post.Publish_Date,
				Close_Date:     post.Close_Date,
				Category:       post.Category.Name,
				Country:        post.Post.Country.Name,
			})
		}
		// ส่งข้อมูลกลับไปในรูปแบบ JSON
		return ctx.Status(200).JSON(postsResponse)
	}
}

func CreatePost(ctx fiber.Ctx) error {
	post := new(request.AnnouncePostCreateRequest)
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
	} else {
		// Handle File Image
		filename := ctx.Locals("filenameImage").(*string)
		// Handle File Attach
		filenameAttach := ctx.Locals("filenameAttach").(*string)

		if post.Publish_Date == nil {
			now := time.Now()
			post.Publish_Date = &now
		}

		// เริ่มต้น Transaction
		tx := database.DB.Begin()
		if tx.Error != nil {
			return ctx.Status(404).JSON(fiber.Map{
				"error": "Failed to begin transaction",
			})
		}

		// Create New Post
		newPost := entity.Post{
			Title:        post.Title,
			Description:  post.Description,
			Image:        filename,
			Posts_Type:   post.Posts_Type,
			Publish_Date: post.Publish_Date,
			Country_ID:   post.Country_ID,
		}
		if err := tx.Create(&newPost).Error; err != nil {
			tx.Rollback()
			return ctx.Status(404).JSON(fiber.Map{
				"error": "Failed to create post",
			})
		}
		// ตรวจสอบว่าได้ Posts_ID หลังจากการสร้าง post
		if newPost.Posts_ID == 0 {
			tx.Rollback()
			return ctx.Status(404).JSON(fiber.Map{
				"error": "Failed to create post",
			})
		}

		// Create New Announce Post
		newAnnouncePost := entity.Announce_Post{
			Posts_ID:    newPost.Posts_ID,
			Url:         post.URL,
			Attach_File: filenameAttach,
			Close_Date:  post.Close_Date,
			Category_ID: post.Category_ID,
		}
		if err := tx.Create(&newAnnouncePost).Error; err != nil {
			tx.Rollback()
			return ctx.Status(404).JSON(fiber.Map{
				"error": "Failed to create announce post",
			})
		}

		// ยืนยันการทำงานของ Transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return ctx.Status(404).JSON(fiber.Map{
				"error": "Failed to commit transaction",
			})
		}

		// สร้างตัวแปรแบบ AnnouncePostResponse
		// และกำหนดค่าให้กับตัวแปรนี้
		postResponse := response.AnnouncePostResponseAdd{
			Announce_ID:    newAnnouncePost.Announce_ID,
			Title:          newPost.Title,
			Description:    newPost.Description,
			URL:            newAnnouncePost.Url,
			Attach_File:    newAnnouncePost.Attach_File,
			Image:          newPost.Image,
			Post_Type:      newPost.Posts_Type,
			Published_Date: newPost.Publish_Date,
			Close_Date:     newAnnouncePost.Close_Date,
			Category_ID:    newAnnouncePost.Category_ID,
			Country_ID:     newPost.Country_ID,
		}
		// ส่งข้อมูลกลับไปในรูปแบบ JSON
		return ctx.Status(201).JSON(postResponse)

	}
	return ctx.Status(400).JSON(fiber.Map{
		"error": "Bad Request",
	})
}

// func GetPostByID(ctx fiber.Ctx) error {
// 	postId := ctx.Params("id")
// 	var post entity.Post
// 	err := database.DB.First(&post, "posts_id = ?", postId).Error
// 	if err != nil {
// 		return ctx.Status(404).JSON(fiber.Map{
// 			"error message": err.Error(),
// 		})
// 	}
// 	return ctx.JSON(post)
// }

// func UpdatePost(ctx fiber.Ctx) error {
// 	postRequest := new(request.PostUpdateRequest)
// 	if err := ctx.Bind().Body(postRequest); err != nil {
// 		return ctx.Status(400).JSON(fiber.Map{
// 			"error message": "Bad Request",
// 		})
// 	}
// 	var post entity.Post
// 	postId := ctx.Params("id")
// 	// Check Available Post
// 	err := database.DB.First(&post, "posts_id = ?", postId).Error
// 	if err != nil {
// 		return ctx.Status(404).JSON(fiber.Map{
// 			"error message": "Post not found",
// 		})
// 	} else {
// 		// Update File Image
// 		_, errFile := ctx.FormFile("image")
// 		if errFile != nil {
// 			log.Println("Error File = ", errFile)
// 		} else {
// 			// Remove Old File Image
// 			if post.Image != nil {
// 				errDeleteImage := utils.HandleRemoveFileImage(*post.Image)
// 				if errDeleteImage != nil {
// 					log.Println("Failed to remove image file:", errDeleteImage)
// 				}
// 			}
// 			// Add New File Image
// 			filename := ctx.Locals("filenameImage").(*string)
// 			post.Image = filename
// 		}

// 		// Update File Attach
// 		_, errFileAttach := ctx.FormFile("attach_file")
// 		if errFileAttach != nil {
// 			log.Println("Error File Attach = ", errFileAttach)
// 		} else {
// 			// Remove Old File Attach
// 			if post.Attach_File != nil {
// 				errDeleteAttach := utils.HandleRemoveFileAttach(*post.Attach_File)
// 				if errDeleteAttach != nil {
// 					log.Println("Failed to remove attach file:", errDeleteAttach)
// 				}
// 			}
// 			// Add New File Attach
// 			filenameAttach := ctx.Locals("filenameAttach").(*string)
// 			post.Attach_File = filenameAttach
// 		}
// 		// Update Post
// 		if postRequest.Title != "" {
// 			post.Title = postRequest.Title
// 		}
// 		if postRequest.Description != "" {
// 			post.Description = postRequest.Description
// 		}
// 		post.URL = postRequest.URL
// 		post.Close_Date = postRequest.Close_Date

// 		errUpdate := database.DB.Save(&post).Error
// 		if errUpdate != nil {
// 			return ctx.Status(400).JSON(fiber.Map{
// 				"error message": errUpdate.Error(),
// 			})
// 		}
// 	}
// 	return ctx.Status(200).JSON(post)
// }

// func DeletePost(ctx fiber.Ctx) error {
// 	postId := ctx.Params("id")
// 	var post entity.Post
// 	// Check Available Post
// 	err := database.DB.Debug().First(&post, "posts_id = ?", postId).Error
// 	if err != nil {
// 		return ctx.Status(404).JSON(fiber.Map{
// 			"message": "post not found",
// 		})
// 	} else {
// 		// Handle File Remove Image
// 		if post.Image != nil {
// 			errDeleteImage := utils.HandleRemoveFileImage(*post.Image)
// 			if errDeleteImage != nil {
// 				log.Println("Failed to remove image file:", errDeleteImage)
// 			}
// 		}

// 		// Handle File Remove Attach
// 		if post.Attach_File != nil {
// 			errDeleteAttach := utils.HandleRemoveFileAttach(*post.Attach_File)
// 			if errDeleteAttach != nil {
// 				log.Println("Failed to remove attach file:", errDeleteAttach)
// 			}
// 		}

// 		errDeletePost := database.DB.Debug().Delete(&post, "posts_id = ?", postId).Error
// 		if errDeletePost != nil {
// 			return ctx.Status(400).JSON(fiber.Map{
// 				"error message": errDeletePost.Error(),
// 			})
// 		}
// 		return ctx.JSON(fiber.Map{
// 			"message": "post deleted",
// 		})
// 	}
// }
