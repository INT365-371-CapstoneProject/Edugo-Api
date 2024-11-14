package handler

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func GetAllAnnouncePost(ctx fiber.Ctx) error {
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

func CreateAnnouncePost(ctx fiber.Ctx) error {
	post := new(request.AnnouncePostCreateRequest)
	if err := ctx.Bind().Body(post); err != nil {
		// clear temp file
		utils.ClearTempFiles()
		// create temp folder
		utils.CreateTempFolder()
		return err
	}
	// Validate Request
	validate := validator.New()
	if errValidate := validate.Struct(post); errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				utils.ClearTempFiles()
				utils.CreateTempFolder()
				return ctx.Status(400).JSON(fiber.Map{
					"error": "The post type field is required.",
				})
			case "oneof":
				utils.ClearTempFiles()
				utils.CreateTempFolder()
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
		// ย้ายไฟล์จาก temp ไปยัง public
		utils.RemoveTempToPublic()
		// ส่งข้อมูลกลับไปในรูปแบบ JSON
		return ctx.Status(201).JSON(postResponse)

	}
	return ctx.Status(400).JSON(fiber.Map{
		"error": "Bad Request",
	})
}

func GetAnnouncePostByID(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post []entity.Announce_Post
	result := database.DB.Where("announce_id = ?", postId).First(&post)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error message": result.Error.Error(),
		})
	} else {
		utils.GetCategoryName(post)
		utils.GetPostByAnnounceID(post)

		var postsResponse []response.AnnouncePostResponse
		for _, post := range post {
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

func UpdateAnnouncePost(ctx fiber.Ctx) error {
	// Bind the update request data
	postRequest := new(request.AnnouncePostUpdateRequest)
	if err := ctx.Bind().Body(postRequest); err != nil {
		// clear temp file
		utils.ClearTempFiles()
		// create temp folder
		utils.CreateTempFolder()
		return ctx.Status(400).JSON(fiber.Map{
			"error message": "Invalid request data",
		})
	}

	postId := ctx.Params("id")

	// Find the existing announce post and preload the associated post
	var announcePost entity.Announce_Post
	err := database.DB.Preload("Post").Where("announce_id = ?", postId).First(&announcePost).Error
	if err != nil {
		// clear temp file
		utils.ClearTempFiles()
		// create temp folder
		utils.CreateTempFolder()
		return ctx.Status(404).JSON(fiber.Map{
			"error message": "Post not found",
		})
	}

	// Begin a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error": "Failed to begin transaction",
		})
	}

	// Update fields in Post table based on request data
	if postRequest.Title != "" {
		announcePost.Post.Title = postRequest.Title
	}
	if postRequest.Description != "" {
		announcePost.Post.Description = postRequest.Description
	}
	if postRequest.Publish_Date != nil {
		announcePost.Post.Publish_Date = postRequest.Publish_Date
	}
	if postRequest.Country_ID != 0 {
		announcePost.Post.Country_ID = postRequest.Country_ID
	}
	if postRequest.Category_ID != 0 {
		announcePost.Category_ID = postRequest.Category_ID
	}

	// Update File Image if provided
	if _, errFile := ctx.FormFile("image"); errFile == nil {
		// Remove old file if exists
		if announcePost.Post.Image != nil {
			if err := utils.HandleRemoveFileImage(*announcePost.Post.Image); err != nil {
				log.Println("Failed to remove old image file:", err)
			}
		}
		// Set new file
		filename := ctx.Locals("filenameImage").(*string)
		announcePost.Post.Image = filename
	}

	// Update File Attach if provided
	if _, errFileAttach := ctx.FormFile("attach_file"); errFileAttach == nil {
		// Remove old attach file if exists
		if announcePost.Attach_File != nil {
			if err := utils.HandleRemoveFileAttach(*announcePost.Attach_File); err != nil {
				log.Println("Failed to remove old attach file:", err)
			}
		}
		// Set new attach file
		filenameAttach := ctx.Locals("filenameAttach").(*string)
		announcePost.Attach_File = filenameAttach
	}

	// Save updated Post record
	if err := tx.Save(&announcePost.Post).Error; err != nil {
		tx.Rollback()
		return ctx.Status(404).JSON(fiber.Map{
			"error message": "Failed to update post details",
		})
	}

	// Save updated Announce_Post record
	if err := tx.Save(&announcePost).Error; err != nil {
		tx.Rollback()
		return ctx.Status(404).JSON(fiber.Map{
			"error message": "Failed to update announce post details",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return ctx.Status(404).JSON(fiber.Map{
			"error message": "Failed to commit transaction",
		})
	}

	// Construct response data
	postResponse := response.AnnouncePostResponseAdd{
		Announce_ID:    announcePost.Announce_ID,
		Title:          announcePost.Post.Title,
		Description:    announcePost.Post.Description,
		URL:            announcePost.Url,
		Attach_File:    announcePost.Attach_File,
		Image:          announcePost.Post.Image,
		Post_Type:      announcePost.Post.Posts_Type,
		Published_Date: announcePost.Post.Publish_Date,
		Close_Date:     announcePost.Close_Date,
		Category_ID:    announcePost.Category_ID,
		Country_ID:     announcePost.Post.Country_ID,
	}

	// Move files from temp to public
	utils.RemoveTempToPublic()
	// Return the updated response
	return ctx.Status(200).JSON(postResponse)
}

func DeleteAnnouncePost(ctx fiber.Ctx) error {
	// รับค่า ID ของ Post จากพารามิเตอร์
	postId := ctx.Params("id")

	// ค้นหาและดึงข้อมูล Announce_Post
	var announcePost entity.Announce_Post
	err := database.DB.First(&announcePost, "posts_id = ?", postId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "announce post not found",
		})
	}

	// ค้นหาและดึงข้อมูล Post ที่เกี่ยวข้อง
	var post entity.Post
	err = database.DB.First(&post, "posts_id = ?", postId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "post not found",
		})
	}

	// จัดการลบไฟล์รูปภาพ (ถ้ามี)
	if post.Image != nil {
		errDeleteImage := utils.HandleRemoveFileImage(*post.Image)
		if errDeleteImage != nil {
			log.Println("Failed to remove image file:", errDeleteImage)
		}
	}

	// จัดการลบไฟล์แนบ (ถ้ามี)
	if announcePost.Attach_File != nil {
		errDeleteAttach := utils.HandleRemoveFileAttach(*announcePost.Attach_File)
		if errDeleteAttach != nil {
			log.Println("Failed to remove attach file:", errDeleteAttach)
		}
	}

	// เริ่มต้น Transaction เพื่อให้แน่ใจว่าการลบสำเร็จหรือยกเลิกทั้งหมดหากเกิดปัญหา
	tx := database.DB.Begin()

	// ลบข้อมูล Announce_Post
	if err := tx.Delete(&announcePost, "posts_id = ?", postId).Error; err != nil {
		tx.Rollback()
		return ctx.Status(400).JSON(fiber.Map{
			"error message": "failed to delete announce post",
		})
	}

	// ลบข้อมูล Post
	if err := tx.Delete(&post, "posts_id = ?", postId).Error; err != nil {
		tx.Rollback()
		return ctx.Status(400).JSON(fiber.Map{
			"error message": "failed to delete post",
		})
	}

	// ยืนยันการทำงานของ Transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return ctx.Status(400).JSON(fiber.Map{
			"error message": "failed to commit transaction",
		})
	}

	// ส่งข้อความตอบกลับ
	return ctx.JSON(fiber.Map{
		"message": "post and announce post deleted",
	})
}
