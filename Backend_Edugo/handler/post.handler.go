package handler

import (
	"log"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
)

// ตัวแปรสำหรับการแปลภาษาและตรวจสอบความถูกต้อง
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

// กำหนดค่าเริ่มต้นสำหรับการตรวจสอบและแปลภาษา
func init() {
	enLocale := en.New()
	uni = ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")
	validate = validator.New()
	enTranslations.RegisterDefaultTranslations(validate, trans)
}

// ฟังก์ชันสำหรับจัดการข้อผิดพลาด
func handleError(ctx fiber.Ctx, statusCode int, message string) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
}

// GetAllPost - ดึงข้อมูลโพสต์ทั้งหมดที่เป็นประเภท Subject
func GetAllPost(ctx fiber.Ctx) error {
	var posts []entity.Post
	result := database.DB.Where("posts_type = ?", "Subject").Find(&posts)
	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	var postsResponse []response.PostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, response.PostResponse{
			Post_ID:      post.Posts_ID,
			Description:  post.Description,
			Image:        post.Image,
			Publish_Date: post.Publish_Date,
			Posts_Type:   post.Posts_Type,
		})
	}
	return ctx.Status(200).JSON(postsResponse)
}

// GetPostByID - ดึงข้อมูลโพสต์ตาม ID ที่ระบุ
func GetPostByID(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post entity.Post
	
	result := database.DB.Where("posts_id = ? AND posts_type = ?", postId, "Subject").First(&post)
	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	postResponse := response.PostResponse{
		Post_ID:      post.Posts_ID,
		Description:  post.Description,
		Image:        post.Image,
		Publish_Date: post.Publish_Date,
		Posts_Type:   post.Posts_Type,
	}
	return ctx.Status(200).JSON(postResponse)
}

// CreatePost - สร้างโพสต์ใหม่
func CreatePost(ctx fiber.Ctx) error {
	post := new(request.PostCreateRequest)
	if err := ctx.Bind().Body(post); err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Invalid request data")
	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := validate.Struct(post); err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, err.(validator.ValidationErrors)[0].Translate(trans))
	}

	// จัดการไฟล์รูปภาพ
	filename := ctx.Locals("filenameImage").(*string)
	
	// กำหนดวันที่เผยแพร่
	if post.Publish_Date == nil {
		now := time.Now().UTC()
		post.Publish_Date = &now
	} else {
		utcTime := post.Publish_Date.UTC()
		post.Publish_Date = &utcTime
	}

	// เริ่ม Transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to begin transaction")
	}

	// Create New Post
	newPost := entity.Post{
		Description:  post.Description,
		Image:        filename,
		Posts_Type:   post.Posts_Type,
		Publish_Date: post.Publish_Date,
	}
	if err := tx.Create(&newPost).Error; err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		tx.Rollback()
		return handleError(ctx, 400, "Failed to create post")
	}

	// ยืนยันการทำงานของ Transaction
	if err := tx.Commit().Error; err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		tx.Rollback()
		return handleError(ctx, 400, "Failed to commit transaction")
	}

	// สร้างตัวแปรแบบ PostResponse
	// และกำหนดค่าให้กับตัวแปรนี้
	postResponse := response.PostResponseAdd{
		Post_ID:      newPost.Posts_ID,
		Description:  newPost.Description,
		Image:        newPost.Image,
		Publish_Date: newPost.Publish_Date,
		Posts_Type:   newPost.Posts_Type,
	}
	// ย้ายไฟล์จาก temp ไปยัง public
	utils.RemoveTempToPublic()
	// ส่งข้อมูลกลับไปในรูปแบบ JSON
	return ctx.Status(201).JSON(postResponse)
}

// UpdatePost - อัปเดตโพสต์
func UpdatePost(ctx fiber.Ctx) error {
	postRequest := new(request.PostUpdateRequest)
	if err := ctx.Bind().Body(postRequest); err != nil {
		// clear temp file
		utils.ClearTempFiles()
		// create temp folder
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Invalid request data")
	}
	postId := ctx.Params("id")

	var post entity.Post
	err := database.DB.Where("posts_id = ? AND posts_type = ?", postId, "Subject").First(&post).Error
	if err != nil {
		// clear temp file
		utils.ClearTempFiles()
		// create temp folder
		utils.CreateTempFolder()
		return handleError(ctx, 404, "Post not found")
	}

	// Validate Request
	if errValidate := validate.Struct(postRequest); errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			utils.ClearTempFiles()
			utils.CreateTempFolder()
			return handleError(ctx, 400, err.Translate(trans))
		}
	}

	// Begin a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to begin transaction")
	}

	// Update fields in Post table based on request data
	if postRequest.Description != "" {
		post.Description = postRequest.Description
	}
	if postRequest.Publish_Date != nil {
		utcTime := postRequest.Publish_Date.UTC()
		post.Publish_Date = &utcTime
	}

	// Update File Image if provided
	if _, errFile := ctx.FormFile("image"); errFile == nil {
		// Remove old file if exists
		if post.Image != nil {
			if err := utils.HandleRemoveFileImage(*post.Image); err != nil {
				log.Println("Failed to remove old image file:", err)
			}
		}
		// Set new file
		filename := ctx.Locals("filenameImage").(*string)
		post.Image = filename
	}

	// Save updated Post record
	if err := tx.Save(&post).Error; err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		tx.Rollback()
		return handleError(ctx, 400, "Failed to update post details")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		tx.Rollback()
		return handleError(ctx, 400, "Failed to commit transaction")
	}

	// Construct response data
	postResponse := response.PostResponseAdd{
		Post_ID:      post.Posts_ID,
		Description:  post.Description,
		Image:        post.Image,
		Publish_Date: post.Publish_Date,
		Posts_Type:   post.Posts_Type,
	}

	// Move files from temp to public
	utils.RemoveTempToPublic()
	// Return the updated response
	return ctx.Status(200).JSON(postResponse)
}

// DeletePost - ลบโพสต์
func DeletePost(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post entity.Post
	err := database.DB.Where("posts_id = ? AND posts_type = ?", postId, "Subject").First(&post).Error
	if err != nil {
		return handleError(ctx, 404, "post not found")
	}

	if post.Image != nil {
		errDeleteImage := utils.HandleRemoveFileImage(*post.Image)
		if errDeleteImage != nil {
			log.Println("Failed to remove image file:", errDeleteImage)
		}
	}

	err = database.DB.Delete(&post).Error
	if err != nil {
		return handleError(ctx, 400, "failed to delete post")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "post deleted successfully",
	})
}

// GetAllAnnouncePost - ดึงข้อมูลประกาศทั้งหมด
func GetAllAnnouncePost(ctx fiber.Ctx) error {
	var posts []entity.Announce_Post
	result := database.DB.Preload("Post").Preload("Category").Preload("Country").Find(&posts)
	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	var postsResponse []response.AnnouncePostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, response.AnnouncePostResponse{
			Announce_ID:  post.Announce_ID,
			Title:        post.Title,
			Description:  post.Post.Description,
			URL:          post.Url,
			Attach_File:  post.Attach_File,
			Image:        post.Post.Image,
			Posts_Type:   post.Post.Posts_Type,
			Publish_Date: post.Post.Publish_Date,
			Close_Date:   post.Close_Date,
			Category:     post.Category.Name,
			Country:      post.Country.Name,
		})
	}
	return ctx.Status(200).JSON(postsResponse)
}

// GetAnnouncePostByID - ดึงข้อมูลประกาศตาม ID ที่ระบุ
func GetAnnouncePostByID(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post []entity.Announce_Post
	result := database.DB.Where("announce_id = ?", postId).Preload("Category").Preload("Country").Preload("Post").First(&post)
	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	postsResponse := response.AnnouncePostResponse{
		Announce_ID:  post[0].Announce_ID,
		Title:        post[0].Title,
		Description:  post[0].Post.Description,
		URL:          post[0].Url,
		Attach_File:  post[0].Attach_File,
		Image:        post[0].Post.Image,
		Posts_Type:   post[0].Post.Posts_Type,
		Publish_Date: post[0].Post.Publish_Date,
		Close_Date:   post[0].Close_Date,
		Category:     post[0].Category.Name,
		Country:      post[0].Country.Name,
	}
	return ctx.Status(200).JSON(postsResponse)
}

// CreateAnnouncePost - สร้างประกาศใหม่
func CreateAnnouncePost(ctx fiber.Ctx) error {
	post := new(request.AnnouncePostCreateRequest)
	if err := ctx.Bind().Body(post); err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Invalid request data")
	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := validate.Struct(post); err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, err.(validator.ValidationErrors)[0].Translate(trans))
	}

	// จัดการไฟล์รูปภาพ
	filename := ctx.Locals("filenameImage").(*string)
	// จัดการไฟล์แนบ
	filenameAttach := ctx.Locals("filenameAttach").(*string)

	// กำหนดวันที่เผยแพร่
	if post.Publish_Date == nil {
		now := time.Now().UTC()
		post.Publish_Date = &now
	} else {
		utcTime := post.Publish_Date.UTC()
		post.Publish_Date = &utcTime
	}

	// กำหนดวันที่ปิดประกาศ
	if post.Close_Date != nil {
		utcTime := post.Close_Date.UTC()
		if utcTime.Before(*post.Publish_Date) {
			utils.ClearTempFiles()
			utils.CreateTempFolder()
			return handleError(ctx, 400, "Close date cannot be before publish date")
		}
		post.Close_Date = &utcTime
	}

	// เริ่ม Transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to begin transaction")
	}

	// Create New Post
	newPost := entity.Post{
		Description:  post.Description,
		Image:        filename,
		Posts_Type:   post.Posts_Type,
		Publish_Date: post.Publish_Date,
	}
	if err := tx.Create(&newPost).Error; err != nil {
		tx.Rollback()
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to create post")
	}

	// ตรวจสอบว่าได้ Posts_ID หลังจากการสร้าง post
	if newPost.Posts_ID == 0 {
		tx.Rollback()
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to create post")
	}

	// Create New Announce Post
	newAnnouncePost := entity.Announce_Post{
		Title:       post.Title,
		Posts_ID:    newPost.Posts_ID,
		Url:         post.URL,
		Attach_File: filenameAttach,
		Close_Date:  post.Close_Date,
		Category_ID: post.Category_ID,
		Country_ID:  post.Country_ID,
	}
	if err := tx.Create(&newAnnouncePost).Error; err != nil {
		tx.Rollback()
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to create announce post")
	}

	// ยืนยันการทำงานของ Transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to commit transaction")
	}

	// สร้างตัวแปรแบบ AnnouncePostResponse
	// และกำหนดค่าให้กับตัวแปรนี้
	postResponse := response.AnnouncePostResponseAdd{
		Announce_ID:  newAnnouncePost.Announce_ID,
		Title:        newAnnouncePost.Title,
		Description:  newPost.Description,
		URL:          newAnnouncePost.Url,
		Attach_File:  newAnnouncePost.Attach_File,
		Image:        newPost.Image,
		Posts_Type:   newPost.Posts_Type,
		Publish_Date: newPost.Publish_Date,
		Close_Date:   newAnnouncePost.Close_Date,
		Category_ID:  newAnnouncePost.Category_ID,
		Country_ID:   newAnnouncePost.Country_ID,
	}
	// ย้ายไฟล์จาก temp ไปยัง public
	utils.RemoveTempToPublic()
	// ส่งข้อมูลกลับไปในรูปแบบ JSON
	return ctx.Status(201).JSON(postResponse)
}

// UpdateAnnouncePost - อัปเดตประกาศ
func UpdateAnnouncePost(ctx fiber.Ctx) error {
	// Bind the update request data
	postRequest := new(request.AnnouncePostUpdateRequest)
	if err := ctx.Bind().Body(postRequest); err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Invalid request data")
	}

	postId := ctx.Params("id")

	// Find the existing announce post and preload the associated post
	var announcePost entity.Announce_Post
	err := database.DB.Preload("Post").Where("announce_id = ?", postId).First(&announcePost).Error
	if err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 404, "Post not found")
	}

	// Validate Request
	if err := validate.Struct(postRequest); err != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, err.(validator.ValidationErrors)[0].Translate(trans))
	}

	// Begin a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to begin transaction")
	}

	// Update fields in Post table based on request data
	if postRequest.Title != "" {
		announcePost.Title = postRequest.Title
	}
	if postRequest.Description != "" {
		announcePost.Post.Description = postRequest.Description
	}
	if postRequest.URL != nil {
		announcePost.Url = postRequest.URL
	}
	if postRequest.Publish_Date != nil {
		utcTime := postRequest.Publish_Date.UTC()
		announcePost.Post.Publish_Date = &utcTime
	}
	if postRequest.Close_Date != nil {
		utcTime := postRequest.Close_Date.UTC()
		if utcTime.Before(*announcePost.Post.Publish_Date) {
			utils.ClearTempFiles()
			utils.CreateTempFolder()
			return handleError(ctx, 400, "Close date cannot be before publish date")
		}
		announcePost.Close_Date = &utcTime
	}
	if postRequest.Country_ID != 0 {
		announcePost.Country_ID = postRequest.Country_ID
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
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to update post details")
	}

	// Save updated Announce_Post record
	if err := tx.Save(&announcePost).Error; err != nil {
		tx.Rollback()
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to update announce post details")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.ClearTempFiles()
		utils.CreateTempFolder()
		return handleError(ctx, 400, "Failed to commit transaction")
	}

	// Construct response data
	postResponse := response.AnnouncePostResponseAdd{
		Announce_ID:  announcePost.Announce_ID,
		Title:        announcePost.Title,
		Description:  announcePost.Post.Description,
		URL:          announcePost.Url,
		Attach_File:  announcePost.Attach_File,
		Image:        announcePost.Post.Image,
		Posts_Type:   announcePost.Post.Posts_Type,
		Publish_Date: announcePost.Post.Publish_Date,
		Close_Date:   announcePost.Close_Date,
		Category_ID:  announcePost.Category_ID,
		Country_ID:   announcePost.Country_ID,
	}

	// Move files from temp to public
	utils.RemoveTempToPublic()
	// Return the updated response
	return ctx.Status(200).JSON(postResponse)
}

// DeleteAnnouncePost - ลบประกาศ
func DeleteAnnouncePost(ctx fiber.Ctx) error {
	// รับค่า ID ของ Post จากพารามิเตอร์
	postId := ctx.Params("id")

	// ค้นหาและดึงข้อมูล Announce_Post
	var announcePost entity.Announce_Post
	err := database.DB.Preload("Post").First(&announcePost, "announce_id = ?", postId).Error
	if err != nil {
		return handleError(ctx, 404, "announce post not found")
	}

	// ค้นหาและดึงข้อมูล Post ที่เกี่ยวข้อง
	var post entity.Post
	err = database.DB.First(&post, "posts_id = ?", announcePost.Posts_ID).Error
	if err != nil {
		return handleError(ctx, 404, "post not found")
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
	if err := tx.Delete(&announcePost, "announce_id = ?", postId).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 400, "failed to delete announce post")
	}

	// ลบข้อมูล Post
	if err := tx.Delete(&post, "posts_id = ?", announcePost.Posts_ID).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 400, "failed to delete post")
	}

	// ยืนยันการทำงานของ Transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 400, "failed to commit transaction")
	}

	// ส่งข้อความตอบกลับ
	return ctx.Status(200).JSON(fiber.Map{
		"message": "post and announce post deleted",
	})
}