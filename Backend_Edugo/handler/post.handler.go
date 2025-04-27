package handler

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/middleware"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
	"gorm.io/gorm" // Add this import
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

	// Register custom translations
	translations := []struct {
		tag         string
		translation string
	}{
		{"required", "{0} is required"},
		{"min", "{0} must be at least {1} characters"},
		{"max", "{0} must be at most {0} characters"},
		{"url", "URL must be a valid URL"},
		{"oneof", "{0} must be {1}"},
	}

	// Register all translations
	for _, t := range translations {
		validate.RegisterTranslation(t.tag, trans, func(ut ut.Translator) error {
			return ut.Add(t.tag, t.translation, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			param := fe.Param()
			if t.tag == "oneof" {
				param = "'" + param + "'"
			}
			t, _ := ut.T(t.tag, fe.Field(), param)
			return t
		})
	}

	// Register custom messages for specific fields
	fieldMessages := map[string]string{
		"Title.required":       "Title is required and must be between 5-100 characters",
		"Description.required": "Description is required and must be between 10-500 characters",
		"Posts_Type.required":  "Posts_Type is required and must be 'Subject' or 'Announce'",
		"Close_Date.required":  "Close_Date is required",
		"Category_ID.required": "Category_ID is required",
		"Country_ID.required":  "Country_ID is required",
	}

	for field, msg := range fieldMessages {
		parts := strings.Split(field, ".")
		validate.RegisterTranslation(parts[1], trans, func(ut ut.Translator) error {
			return ut.Add(field, msg, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			if fe.Field() == parts[0] {
				t, _ := ut.T(field)
				return t
			}
			return fe.Error()
		})
	}

	// Add this validator in the init() function
	validate.RegisterValidation("education_level", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		validLevels := map[string]bool{
			"Undergraduate": true,
			"Master":        true,
			"Doctorate":     true,
		}
		return validLevels[value]
	})
}

// ฟังก์ชันสำหรับจัดการข้อผิดพลาด
func handleError(ctx fiber.Ctx, statusCode int, message string) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
}

// checkRole verifies if the user has the required role(s)
func checkRole(ctx fiber.Ctx, allowedRoles ...string) (string, string, error) {
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)
	username := claims["username"].(string)

	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return role, username, nil
		}
	}
	return "", "", handleError(ctx, fiber.StatusUnauthorized,
		"Unauthorized: requires one of these roles: "+strings.Join(allowedRoles, ", "))
}

// getAccount retrieves account information from username
func getAccount(username string) (*entity.Account, error) {
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// getPaginationParams extracts page and limit from query parameters
func getPaginationParams(ctx fiber.Ctx) (page, limit, offset int) {
	page = 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	limit = 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = limitNum
		}
	}

	offset = (page - 1) * limit
	return page, limit, offset
}

func GetAllPost(ctx fiber.Ctx) error {
	var posts []struct {
		entity.Post
		Fullname string `json:"fullname"`
	}

	// ดึงข้อมูลทั้งหมดพร้อม JOIN กับ accounts
	result := database.DB.Table("posts p").
		Select(`p.*, 
            CASE 
                WHEN pr.company_name IS NOT NULL THEN pr.company_name
                ELSE CONCAT(a.first_name, ' ', a.last_name)
            END AS fullname`).
		Joins("JOIN accounts a ON p.account_id = a.account_id").
		Joins("LEFT JOIN providers pr ON a.account_id = pr.account_id").
		Scan(&posts)

	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	// แปลงข้อมูลเป็น response
	var postsResponse []response.PostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, response.PostResponse{
			Post_ID:      post.Posts_ID,
			Description:  post.Description,
			Publish_Date: post.Publish_Date,
			Account_ID:   post.Account_ID,
			Fullname:     post.Fullname,
		})
	}

	// ส่งข้อมูลกลับโดยตรง ไม่ใช้ pagination
	return ctx.Status(200).JSON(postsResponse)
}

// GetPostByID - ดึงข้อมูลโพสต์ตาม ID ที่ระบุ
func GetPostByID(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post entity.Post

	result := database.DB.Where("posts_id = ?", postId).First(&post)
	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	postResponse := response.PostResponse{
		Post_ID:      post.Posts_ID,
		Description:  post.Description,
		Publish_Date: post.Publish_Date,
	}
	return ctx.Status(200).JSON(postResponse)
}

// GetPostImage - ดึงข้อมูลรูปภาพของโพสต์
func GetPostImage(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post entity.Post
	result := database.DB.Where("posts_id = ?", postId).First(&post)
	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	if post.Image == nil {
		return handleError(ctx, 404, "Image not found")
	}

	ctx.Set("Content-Type", "image/jpeg")
	return ctx.Send(post.Image)
}

// GetPostImage - ดึงข้อมูลรูปภาพของโพสต์
func GetPostAvatarByAccountID(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post entity.Post
	result := database.DB.Where("posts_id = ?", postId).First(&post)
	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	var account entity.Account
	if err := database.DB.Select("avatar").First(&account, "account_id = ?", post.Account_ID).Error; err != nil {
		return utils.HandleError(ctx, 404, "Avatar not found")
	}
	// If no avatar is stored
	if len(account.Avatar) == 0 {
		return utils.HandleError(ctx, 404, "No avatar image found")
	}

	// Set content type header for image
	ctx.Set("Content-Type", "image/jpeg") // You might want to store the content type in DB if you support multiple formats

	// Return the image bytes directly
	return ctx.Send(account.Avatar)
}

// CreatePost - สร้างโพสต์ใหม่
func CreatePost(ctx fiber.Ctx) error {
	// เรียกใช้ฟังก์ชัน GetTokenClaims เพื่อดึงข้อมูลจาก JWT
	claims := middleware.GetTokenClaims(ctx)
	username := claims["username"].(string)

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return handleError(ctx, 404, "Account not found")
	}

	post := new(request.PostCreateRequest)
	if err := ctx.Bind().Body(post); err != nil {
		return handleError(ctx, 400, "Invalid request data")
	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := validate.Struct(post); err != nil {
		return handleError(ctx, 400, err.(validator.ValidationErrors)[0].Translate(trans))
	}

	// ใช้ฟังก์ชัน HandleImageUpload แทนการจัดการไฟล์โดยตรง
	if err := utils.HandleImageUpload(ctx, "image"); err != nil {
		return handleError(ctx, 400, "Error handling image upload: "+err.Error())
	}

	// กำหนดวันที่เผยแพร่
	if post.Publish_Date == nil {
		now := time.Now().UTC()
		post.Publish_Date = &now
	} else {
		utcTime := post.Publish_Date.UTC()
		post.Publish_Date = &utcTime
	}

	// Create New Post with account ID from JWT
	newPost := entity.Post{
		Description:  post.Description,
		Publish_Date: post.Publish_Date,
		Account_ID:   account.Account_ID, // ใช้ Account_ID จาก JWT
	}

	// ตรวจสอบว่ามีรูปภาพถูกอัพโหลดหรือไม่
	if imageBytes := ctx.Locals("imageBytes"); imageBytes != nil {
		newPost.Image = imageBytes.([]byte)
	}

	// Begin transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return handleError(ctx, 409, "Failed to begin transaction")
	}

	// Create post
	if err := tx.Create(&newPost).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "Failed to create post")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "Failed to commit transaction")
	}

	postResponse := response.PostResponse{
		Post_ID:      newPost.Posts_ID,
		Description:  newPost.Description,
		Publish_Date: newPost.Publish_Date,
		Account_ID:   newPost.Account_ID,
	}

	return ctx.Status(201).JSON(postResponse)
}

// UpdatePost - อัปเดตโพสต์
func UpdatePost(ctx fiber.Ctx) error {
	// เรียกใช้ฟังก์ชัน GetTokenClaims เพื่อดึงข้อมูลจาก JWT
	claims := middleware.GetTokenClaims(ctx)
	username := claims["username"].(string)

	postRequest := new(request.PostUpdateRequest)
	if err := ctx.Bind().Body(postRequest); err != nil {
		return handleError(ctx, 400, "Invalid request data")
	}
	postId := ctx.Params("id")

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	var post entity.Post
	err := database.DB.Where("posts_id = ? AND account_id = ?",
		postId, account.Account_ID).First(&post).Error
	if err != nil {
		return handleError(ctx, 403, "Forbidden")
	}

	// Validate Request
	if errValidate := validate.Struct(postRequest); errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			return handleError(ctx, 400, err.Translate(trans))
		}
	}

	// Handle image upload if provided
	if err := utils.HandleImageUpload(ctx, "image"); err != nil {
		return handleError(ctx, 400, "Error handling image upload: "+err.Error())
	}

	// Begin a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return handleError(ctx, 409, "Failed to begin transaction")
	}

	// Update fields in Post table based on request data
	if postRequest.Description != "" {
		post.Description = postRequest.Description
	}
	if postRequest.Publish_Date != nil {
		utcTime := postRequest.Publish_Date.UTC()
		post.Publish_Date = &utcTime
	}

	// Update image if new image was uploaded
	if imageBytes := ctx.Locals("imageBytes"); imageBytes != nil {
		// If there's a new image, update it
		post.Image = imageBytes.([]byte)
	}

	// Save updated Post record
	if err := tx.Save(&post).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "Failed to update post details")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "Failed to commit transaction")
	}

	// Construct response data
	postResponse := response.PostResponseAdd{
		Post_ID:      post.Posts_ID,
		Description:  post.Description,
		Publish_Date: post.Publish_Date,
		Account_ID:   post.Account_ID,
	}
	// Return the updated response
	return ctx.Status(200).JSON(postResponse)
}

// DeletePost - ลบโพสต์
func DeletePost(ctx fiber.Ctx) error {
	// เรียกใช้ฟังก์ชัน GetTokenClaims เพื่อดึงข้อมูลจาก JWT
	claims := middleware.GetTokenClaims(ctx)
	username := claims["username"].(string)

	postId := ctx.Params("id")

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	var post entity.Post
	err := database.DB.Where("posts_id = ? AND account_id = ?",
		postId, account.Account_ID).First(&post).Error
	if err != nil {
		return handleError(ctx, 403, "Forbidden")
	}
	err = database.DB.Delete(&post).Error
	if err != nil {
		return handleError(ctx, 400, "Failed to delete post")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}

// GetAllAnnouncePost - ดึงข้อมูลประกาศทั้งหมด
func GetAllAnnouncePostForProvider(ctx fiber.Ctx) error {
	// Check role authorization
	_, username, err := checkRole(ctx, "provider")
	if err != nil {
		return err
	}

	// Get account
	account, err := getAccount(username)
	if err != nil {
		return handleError(ctx, 404, "Account not found")
	}

	// เพิ่มการหา provider
	var provider entity.Provider
	if err := database.DB.Where("account_id = ?", account.Account_ID).First(&provider).Error; err != nil {
		return handleError(ctx, 404, "Provider not found")
	}

	// Get pagination parameters
	page, limit, offset := getPaginationParams(ctx)

	var posts []entity.Announce_Post
	var total int64

	// Count total records
	database.DB.Model(&entity.Announce_Post{}).
		Where("provider_id = ?", provider.Provider_ID).
		Count(&total)

	// Get paginated data
	result := database.DB.
		Preload("Category").
		Preload("Country").
		Where("provider_id = ?", provider.Provider_ID).
		Offset(offset).
		Limit(limit).
		Find(&posts)

	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	// Transform to response format
	postsResponse := make([]response.AnnouncePostResponse, len(posts))
	for i, post := range posts {
		postsResponse[i] = response.AnnouncePostResponse{
			Announce_ID:     post.Announce_ID,
			Title:           post.Title,
			Description:     post.Description,
			URL:             post.Url,
			Attach_Name:     post.Attach_Name,
			Publish_Date:    post.Publish_Date,
			Close_Date:      post.Close_Date,
			Category:        post.Category.Name,
			Country:         post.Country.Name,
			Education_Level: post.Education_Level,
			Provider_ID:     post.Provider_ID, // เปลี่ยนจาก Post_ID เป็น Provider_ID
		}
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return ctx.Status(200).JSON(response.PaginatedAnnouncePostResponse{
		Data:     postsResponse,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
		PerPage:  limit,
	})
}

// GetAnnouncePostByID - ดึงข้อมูลประกาศตาม ID ที่ระบุ
func GetAnnouncePostByIDForProvider(ctx fiber.Ctx) error {
	// ตรวจสอบ role จาก JWT
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)
	username := claims["username"].(string)

	// ตรวจสอบว่าเป็น provider เท่านั้น
	if role != "provider" {
		return handleError(ctx, fiber.StatusUnauthorized, "Unauthorized: requires provider role")
	}

	postId := ctx.Params("id")

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	// เพิ่มการหา provider
	var provider entity.Provider
	if err := database.DB.Where("account_id = ?", account.Account_ID).First(&provider).Error; err != nil {
		return handleError(ctx, 404, "Provider not found")
	}

	var post []entity.Announce_Post
	result := database.DB.
		Where("announce_posts.announce_id = ? AND provider_id = ?", postId, provider.Provider_ID).
		Preload("Category").
		Preload("Country").
		First(&post)
	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	postsResponse := response.AnnouncePostResponse{
		Announce_ID:     post[0].Announce_ID,
		Title:           post[0].Title,
		Description:     post[0].Description,
		URL:             post[0].Url,
		Attach_Name:     post[0].Attach_Name,
		Publish_Date:    post[0].Publish_Date,
		Close_Date:      post[0].Close_Date,
		Category:        post[0].Category.Name,
		Country:         post[0].Country.Name,
		Education_Level: post[0].Education_Level,
		Provider_ID:     post[0].Provider_ID, // เปลี่ยนจาก Post_ID เป็น Provider_ID
	}
	return ctx.Status(200).JSON(postsResponse)
}

// GetAnnouncePostAttach - ดึงข้อมูลไฟล์แนบของประกาศ
func GetAnnouncePostAttach(ctx fiber.Ctx) error {
	postId := ctx.Params("id")
	var post entity.Announce_Post
	result := database.DB.Where("announce_id = ?", postId).First(&post)
	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	if post.Attach_File == nil {
		return handleError(ctx, 404, "Attachment not found")
	}

	// ตั้งค่า header สำหรับแสดงผล PDF ในเบราว์เซอร์
	ctx.Set("Content-Type", "application/pdf")
	ctx.Set("Content-Disposition", "inline; filename="+*post.Attach_Name)

	// เพิ่ม cache control เพื่อป้องกันการแคช
	ctx.Set("Cache-Control", "no-store, no-cache, must-revalidate")
	ctx.Set("Pragma", "no-cache")
	ctx.Set("Expires", "0")

	return ctx.Send(post.Attach_File)
}

// GetAnnounceImage - ดึงข้อมูลรูปภาพของประกาศ
func GetAnnounceImage(ctx fiber.Ctx) error {
	postId := ctx.Params("id")

	var announcePost entity.Announce_Post
	result := database.DB.
		Where("announce_posts.announce_id = ?", postId).
		First(&announcePost)

	if result.Error != nil {
		return handleError(ctx, 404, "Announcement not found")
	}

	if announcePost.Image == nil {
		return handleError(ctx, 404, "Image not found")
	}

	ctx.Set("Content-Type", "image/jpeg")
	return ctx.Send(announcePost.Image)
}

// CreateAnnouncePost - สร้างประกาศใหม่
func CreateAnnouncePostForProvider(ctx fiber.Ctx) error {
	// ตรวจสอบ role จาก JWT
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)
	username := claims["username"].(string)

	// ตรวจสอบว่าเป็น provider เท่านั้น
	if role != "provider" {
		return handleError(ctx, fiber.StatusUnauthorized, "Unauthorized: requires provider role")
	}

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return handleError(ctx, 404, "Account not found")
	}

	// เพิ่มการหา provider
	var provider entity.Provider
	if err := database.DB.Where("account_id = ?", account.Account_ID).First(&provider).Error; err != nil {
		return handleError(ctx, 404, "Provider not found")
	}

	post := new(request.AnnouncePostCreateRequest)
	if err := ctx.Bind().Body(post); err != nil {
		return handleError(ctx, 400, "Invalid request data")
	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := validate.Struct(post); err != nil {
		return handleError(ctx, 400, err.(validator.ValidationErrors)[0].Translate(trans))
	}

	// จัดการไฟล์รูปภาพและไฟล์แนบ
	if err := utils.HandleImageUpload(ctx, "image"); err != nil {
		return handleError(ctx, 400, "Error handling image upload: "+err.Error())
	}

	if err := utils.HandleAttachUpload(ctx, "attach_file"); err != nil {
		return handleError(ctx, 400, "Error handling attachment upload: "+err.Error())
	}

	// เก็บชื่อไฟล์แนบ (ถ้ามี)
	var attachFileName *string
	if attachFile, err := ctx.FormFile("attach_file"); err == nil && attachFile != nil {
		fileName := attachFile.Filename
		attachFileName = &fileName
	}

	// กำหนดวันที่
	if post.Publish_Date == nil {
		now := time.Now().UTC()
		post.Publish_Date = &now
	} else {
		utcTime := post.Publish_Date.UTC()
		post.Publish_Date = &utcTime
	}

	if post.Close_Date != nil {
		utcTime := post.Close_Date.UTC()
		if utcTime.Before(*post.Publish_Date) {
			return handleError(ctx, 400, "Close date cannot be before publish date")
		}
		post.Close_Date = &utcTime
	}

	// Begin transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return handleError(ctx, 409, "Failed to begin transaction")
	}

	// Create New Announce Post
	newAnnouncePost := entity.Announce_Post{
		Title:           post.Title,
		Description:     post.Description,
		Url:             post.URL,
		Attach_Name:     attachFileName,
		Publish_Date:    post.Publish_Date,
		Close_Date:      post.Close_Date,
		Category_ID:     post.Category_ID,
		Country_ID:      post.Country_ID,
		Education_Level: post.Education_Level,
		Provider_ID:     provider.Provider_ID,
	}

	// เพิ่มรูปภาพถ้ามีการอัพโหลด
	if imageBytes := ctx.Locals("imageBytes"); imageBytes != nil {
		newAnnouncePost.Image = imageBytes.([]byte)
	}

	// เพิ่มไฟล์แนบถ้ามีการอัพโหลด
	if attachBytes := ctx.Locals("attachBytes"); attachBytes != nil {
		newAnnouncePost.Attach_File = attachBytes.([]byte)
	}

	// Save Announce Post
	if err := tx.Create(&newAnnouncePost).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "Failed to create announce post")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "Failed to commit transaction")
	}

	// Create response
	postResponse := response.AnnouncePostResponseAdd{
		Announce_ID:     newAnnouncePost.Announce_ID,
		Title:           newAnnouncePost.Title,
		Description:     newAnnouncePost.Description,
		URL:             newAnnouncePost.Url,
		Attach_Name:     newAnnouncePost.Attach_Name,
		Publish_Date:    newAnnouncePost.Publish_Date,
		Close_Date:      newAnnouncePost.Close_Date,
		Category_ID:     newAnnouncePost.Category_ID,
		Country_ID:      newAnnouncePost.Country_ID,
		Education_Level: newAnnouncePost.Education_Level,
		Provider_ID:     newAnnouncePost.Provider_ID,
	}

	return ctx.Status(201).JSON(postResponse)
}

// UpdateAnnouncePost - อัปเดตประกาศ
func UpdateAnnouncePostForProvider(ctx fiber.Ctx) error {
	// ตรวจสอบ role จาก JWT
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)
	username := claims["username"].(string)

	// ตรวจสอบว่าเป็น provider เท่านั้น
	if role != "provider" {
		return handleError(ctx, fiber.StatusUnauthorized, "Unauthorized: requires provider role")
	}

	postId := ctx.Params("id")

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return handleError(ctx, 404, "Account not found")
	}

	// เพิ่มการหา provider หลังจากหา account
	var provider entity.Provider
	if err := database.DB.Where("account_id = ?", account.Account_ID).First(&provider).Error; err != nil {
		return handleError(ctx, 404, "Provider not found")
	}

	// Bind the update request data
	postRequest := new(request.AnnouncePostUpdateRequest)
	if err := ctx.Bind().Body(postRequest); err != nil {
		return handleError(ctx, 400, "Invalid request data")
	}

	// แก้ไขส่วนการค้นหา Announce Post
	var announcePost entity.Announce_Post
	err := database.DB.
		Where("announce_id = ? AND provider_id = ?", postId, provider.Provider_ID).
		First(&announcePost).Error
	if err != nil {
		return handleError(ctx, 404, "Post not found")
	}

	// Validate Request
	if err := validate.Struct(postRequest); err != nil {
		return handleError(ctx, 400, err.(validator.ValidationErrors)[0].Translate(trans))
	}
	// Begin a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return handleError(ctx, 409, "Failed to begin transaction")
	}

	// Update fields in Post table based on request data
	if postRequest.Title != "" {
		announcePost.Title = postRequest.Title
	}
	if postRequest.Description != "" {
		announcePost.Description = postRequest.Description
	}
	if postRequest.URL != nil {
		announcePost.Url = postRequest.URL
	}
	if postRequest.Publish_Date != nil {
		utcTime := postRequest.Publish_Date.UTC()
		announcePost.Publish_Date = &utcTime
	}
	if postRequest.Close_Date != nil {
		utcTime := postRequest.Close_Date.UTC()
		if utcTime.Before(*announcePost.Publish_Date) {
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
	if postRequest.Education_Level != "" {
		announcePost.Education_Level = postRequest.Education_Level
	}
	// Update File Image if provided
	if err := utils.HandleImageUpload(ctx, "image"); err == nil {
		if imageBytes := ctx.Locals("imageBytes"); imageBytes != nil {
			announcePost.Image = imageBytes.([]byte)
		}
	}

	// Update File Attach if provided
	if err := utils.HandleAttachUpload(ctx, "attach_file"); err == nil {
		if attachBytes := ctx.Locals("attachBytes"); attachBytes != nil {
			// Update attachment file name
			if attachFile, err := ctx.FormFile("attach_file"); err == nil {
				fileName := attachFile.Filename
				announcePost.Attach_Name = &fileName
			}
			// Update attachment file content
			announcePost.Attach_File = attachBytes.([]byte)
		}
	}

	// Save updated Post record
	// if err := tx.Save(&announcePost.Post).Error; err != nil {
	// 	tx.Rollback()
	// 	return handleError(ctx, 409, "Failed to update post details")
	// }

	// Save updated Announce_Post record
	if err := tx.Save(&announcePost).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "Failed to update announce post details")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "Failed to commit transaction")
	}

	// Construct response data
	postResponse := response.AnnouncePostResponseAdd{
		Announce_ID:     announcePost.Announce_ID,
		Title:           announcePost.Title,
		Description:     announcePost.Description,
		URL:             announcePost.Url,
		Attach_Name:     announcePost.Attach_Name,
		Publish_Date:    announcePost.Publish_Date,
		Close_Date:      announcePost.Close_Date,
		Category_ID:     announcePost.Category_ID,
		Country_ID:      announcePost.Country_ID,
		Provider_ID:     announcePost.Provider_ID,
		Education_Level: announcePost.Education_Level,
	}
	// Return the updated response
	return ctx.Status(200).JSON(postResponse)
}

// DeleteAnnouncePost - ลบประกาศ
func DeleteAnnouncePostForProvider(ctx fiber.Ctx) error {
	// Check role authorization
	_, username, err := checkRole(ctx, "provider")
	if err != nil {
		return err
	}

	// Get account
	account, err := getAccount(username)
	if err != nil {
		return handleError(ctx, 404, "Account not found")
	}

	// เพิ่มการหา provider
	var provider entity.Provider
	if err := database.DB.Where("account_id = ?", account.Account_ID).First(&provider).Error; err != nil {
		return handleError(ctx, 404, "Provider not found")
	}

	postId := ctx.Params("id")

	// Begin transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return handleError(ctx, 409, "failed to begin transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find and verify post ownership
	var announcePost entity.Announce_Post
	if err := tx.Where("announce_id = ? AND provider_id = ?", postId, provider.Provider_ID).
		First(&announcePost).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 404, "announce post not found")
	}

	// Delete announce post directly
	if err := tx.Delete(&announcePost).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "failed to delete announce post")
	}

	if err := tx.Commit().Error; err != nil {
		return handleError(ctx, 409, "failed to commit transaction")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "post deleted successfully",
	})
}

// GetAllAnnouncePostForUser - ดึงข้อมูลประกาศทั้งหมดสำหรับผู้ใช้ทั่วไป
func GetAllAnnouncePostForUser(ctx fiber.Ctx) error {
	// Get page and limit from query parameters
	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	limit := 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = limitNum
		}
	}

	// Calculate offset
	offset := (page - 1) * limit

	var posts []entity.Announce_Post
	var total int64

	// Count total records for active announcements
	database.DB.Model(&entity.Announce_Post{}).
		Where("close_date IS NULL OR close_date > ?", time.Now()).
		Count(&total)

	// Get paginated data with preloaded relations
	result := database.DB.
		Preload("Category").
		Preload("Country").
		Where("close_date IS NULL OR close_date > ?", time.Now()).
		Order("publish_date DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts)

	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	var postsResponse []response.AnnouncePostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, response.AnnouncePostResponse{
			Announce_ID:     post.Announce_ID,
			Title:           post.Title,
			Description:     post.Description,
			URL:             post.Url,
			Attach_Name:     post.Attach_Name,
			Publish_Date:    post.Publish_Date,
			Close_Date:      post.Close_Date,
			Category:        post.Category.Name,
			Country:         post.Country.Name,
			Provider_ID:     post.Provider_ID, // เปลี่ยนจาก Post_ID เป็น Provider_ID
			Education_Level: post.Education_Level,
		})
	}

	// Calculate last page
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	// Create paginated response
	return ctx.Status(200).JSON(response.PaginatedAnnouncePostResponse{
		Data:     postsResponse,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
		PerPage:  limit,
	})
}

// GetAnnouncePostByIDForUser - ดึงข้อมูลประกาศตาม ID สำหรับผู้ใช้ทั่วไป
func GetAnnouncePostByIDForUser(ctx fiber.Ctx) error {
	postId := ctx.Params("id")

	var post entity.Announce_Post
	result := database.DB.
		Preload("Category").
		Preload("Country").
		Where("announce_posts.announce_id = ? AND (announce_posts.close_date IS NULL OR announce_posts.close_date > ?)",
			postId, time.Now()).
		First(&post)

	if result.Error != nil {
		return handleError(ctx, 404, "Announcement not found or no longer available")
	}
	postResponse := response.AnnouncePostResponse{
		Announce_ID:     post.Announce_ID,
		Title:           post.Title,
		Description:     post.Description,
		URL:             post.Url,
		Attach_Name:     post.Attach_Name,
		Publish_Date:    post.Publish_Date,
		Close_Date:      post.Close_Date,
		Category:        post.Category.Name,
		Country:         post.Country.Name,
		Provider_ID:     post.Provider_ID, // เปลี่ยนจาก Post_ID เป็น Provider_ID
		Education_Level: post.Education_Level,
	}

	return ctx.Status(200).JSON(postResponse)
}

func GetAnnouncePostByProviderIDForUser(ctx fiber.Ctx) error {
	providerId := ctx.Params("id")

	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	limit := 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = limitNum
		}
	}

	offset := (page - 1) * limit

	var posts []entity.Announce_Post
	var total int64

	// นับจำนวนทั้งหมดก่อน
	database.DB.Model(&entity.Announce_Post{}).
		Where("provider_id = ? AND (close_date IS NULL OR close_date > ?)", providerId, time.Now()).
		Count(&total)

	// ดึงข้อมูลพร้อมความสัมพันธ์
	result := database.DB.
		Preload("Category").
		Preload("Country").
		Where("provider_id = ? AND (close_date IS NULL OR close_date > ?)", providerId, time.Now()).
		Order("publish_date DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts)

	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	var postsResponse []response.AnnouncePostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, response.AnnouncePostResponse{
			Announce_ID:     post.Announce_ID,
			Title:           post.Title,
			Description:     post.Description,
			URL:             post.Url,
			Attach_Name:     post.Attach_Name,
			Publish_Date:    post.Publish_Date,
			Close_Date:      post.Close_Date,
			Category:        post.Category.Name,
			Country:         post.Country.Name,
			Provider_ID:     post.Provider_ID,
			Education_Level: post.Education_Level,
		})
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return ctx.Status(200).JSON(response.PaginatedAnnouncePostResponse{
		Data:     postsResponse,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
		PerPage:  limit,
	})
}

func GetAnnouncePostByBookmark(ctx fiber.Ctx) error {
	// รับค่าจาก Body
	post := new(request.PostAnnounceBookmarkRequest)
	if err := ctx.Bind().Body(post); err != nil {
		return handleError(ctx, 400, "Invalid request data")
	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := validate.Struct(post); err != nil {
		return handleError(ctx, 400, err.(validator.ValidationErrors)[0].Translate(trans))
	}

	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	limit := 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = limitNum
		}
	}

	// Calculate offset
	offset := (page - 1) * limit
	var total int64

	// Count total records
	database.DB.Model(&entity.Announce_Post{}).
		Where("announce_id IN (?) AND (close_date IS NULL OR close_date > ?)", post.Announce_ID, time.Now()).
		Count(&total)

	// Query ข้อมูลจากฐานข้อมูลโดยใช้ announce_id
	var announcePosts []entity.Announce_Post
	if err := database.DB.Preload("Category").
		Preload("Country").
		Where("announce_id IN (?) AND (close_date IS NULL OR close_date > ?)", post.Announce_ID, time.Now()).
		Offset(offset).
		Limit(limit).
		Find(&announcePosts).Error; err != nil {
		return handleError(ctx, 500, "Database query error")
	}

	// Map ข้อมูลเป็น Response
	var postsResponse []response.AnnouncePostResponse
	for _, post := range announcePosts {
		postsResponse = append(postsResponse, response.AnnouncePostResponse{
			Announce_ID:  post.Announce_ID,
			Title:        post.Title,
			Description:  post.Description,
			URL:          post.Url,
			Attach_Name:  post.Attach_Name,
			Publish_Date: post.Publish_Date,
			Close_Date:   post.Close_Date,
			Category:     post.Category.Name,
			Country:      post.Country.Name,
			Provider_ID:  post.Provider_ID,
		})
	}

	// Calculate last page
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	// Create paginated response
	return ctx.Status(200).JSON(response.PaginatedAnnouncePostResponse{
		Data:     postsResponse,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
		PerPage:  limit,
	})
}

// GetAllAnnouncePostForAdminAndSuperAdmin - ดึงข้อมูลประกาศทั้งหมดสำหรับ Admin และ SuperAdmin
func GetAllAnnouncePostForAdmin(ctx fiber.Ctx) error {
	// ตรวจสอบ role จาก JWT
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)

	// ตรวจสอบว่าเป็น admin หรือ superadmin เท่านั้น
	if role != "admin" && role != "superadmin" {
		return handleError(ctx, fiber.StatusUnauthorized, "Unauthorized: requires admin or superadmin role")
	}

	// ... existing pagination code ...
	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	limit := 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = limitNum
		}
	}

	offset := (page - 1) * limit

	var posts []entity.Announce_Post
	var total int64

	// ... rest of the existing function code ...
	database.DB.Model(&entity.Announce_Post{}).
		Count(&total)

	result := database.DB.
		Preload("Category").
		Preload("Country").
		Order("publish_date DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts)

	if result.Error != nil {
		return handleError(ctx, 404, result.Error.Error())
	}

	// ... rest of response handling ...
	var postsResponse []response.AnnouncePostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, response.AnnouncePostResponse{
			Announce_ID:     post.Announce_ID,
			Title:           post.Title,
			Description:     post.Description,
			URL:             post.Url,
			Attach_Name:     post.Attach_Name,
			Publish_Date:    post.Publish_Date,
			Close_Date:      post.Close_Date,
			Category:        post.Category.Name,
			Country:         post.Country.Name,
			Education_Level: post.Education_Level,
			Provider_ID:     post.Provider_ID, // เปลี่ยนจาก Post_ID เป็น Provider_ID
		})
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return ctx.Status(200).JSON(response.PaginatedAnnouncePostResponse{
		Data:     postsResponse,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
		PerPage:  limit,
	})
}

// GetAnnouncePostByIDForAdminAndSuperAdmin - ดึงข้อมูลประกาศตาม ID สำหรับ Admin และ SuperAdmin
func GetAnnouncePostByIDForAdmin(ctx fiber.Ctx) error {
	// ตรวจสอบ role จาก JWT
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)

	// ตรวจสอบว่าเป็น admin หรือ superadmin เท่านั้น
	if role != "admin" && role != "superadmin" {
		return handleError(ctx, fiber.StatusUnauthorized, "Unauthorized: requires admin or superadmin role")
	}

	postId := ctx.Params("id")

	var post entity.Announce_Post
	result := database.DB.
		Preload("Category").
		Preload("Country").
		Where("announce_posts.announce_id = ?", postId).
		First(&post)

	if result.Error != nil {
		return handleError(ctx, 404, "Announcement not found")
	}

	postResponse := response.AnnouncePostResponse{
		Announce_ID:     post.Announce_ID,
		Title:           post.Title,
		Description:     post.Description,
		URL:             post.Url,
		Attach_Name:     post.Attach_Name,
		Publish_Date:    post.Publish_Date,
		Close_Date:      post.Close_Date,
		Category:        post.Category.Name,
		Country:         post.Country.Name,
		Provider_ID:     post.Provider_ID, // เปลี่ยนจาก Post_ID เป็น Provider_ID
		Education_Level: post.Education_Level,
	}

	return ctx.Status(200).JSON(postResponse)
}

// DeleteAnnouncePostForAdminAndSuperAdmin - ลบประกาศสำหรับ Admin และ SuperAdmin
func DeleteAnnouncePostForAdmin(ctx fiber.Ctx) error {
	// ตรวจสอบ role จาก JWT
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)

	// ตรวจสอบว่าเป็น admin หรือ superadmin เท่านั้น
	if role != "admin" && role != "superadmin" {
		return handleError(ctx, fiber.StatusUnauthorized, "Unauthorized: requires admin or superadmin role")
	}

	postId := ctx.Params("id")

	// Begin transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return handleError(ctx, 409, "failed to begin transaction")
	}

	// Find announce post
	var announcePost entity.Announce_Post
	if err := tx.Where("announce_id = ?", postId).
		First(&announcePost).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 404, "announce post not found")
	}

	// Delete announce post
	if err := tx.Delete(&announcePost).Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "failed to delete post")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return handleError(ctx, 409, "failed to commit transaction")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "post deleted successfully",
	})
}

// DeletePostForAdminAndSuperAdmin - ลบโพสต์สำหรับ Admin และ SuperAdmin
func DeletePostForAdmin(ctx fiber.Ctx) error {
	// ตรวจสอบ role จาก JWT
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)

	// ตรวจสอบว่าเป็น admin หรือ superadmin เท่านั้น
	if role != "admin" && role != "superadmin" {
		return handleError(ctx, fiber.StatusUnauthorized, "Unauthorized: requires admin or superadmin role")
	}

	postId := ctx.Params("id")

	var post entity.Post
	err := database.DB.Where("posts_id = ?", postId).First(&post).Error
	if err != nil {
		return handleError(ctx, 404, "Post not found")
	}

	err = database.DB.Delete(&post).Error
	if err != nil {
		return handleError(ctx, 400, "Failed to delete post")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}

// SearchPosts - แก้ไขการใช้ ILIKE เป็น LIKE with LOWER()
func SearchPosts(ctx fiber.Ctx) error {
	// Get search parameters from query
	search := ctx.Query("search")
	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")
	sortBy := ctx.Query("sort_by", "publish_date") // default sort by publish date
	sortOrder := ctx.Query("sort_order", "desc")   // default descending order

	// Get pagination parameters
	page, limit, offset := getPaginationParams(ctx)

	// Initialize the base query
	query := database.DB.Model(&entity.Post{})

	// Apply search filter if provided
	if search != "" {
		query = query.Where("LOWER(description) LIKE LOWER(?)", "%"+search+"%")
	}

	// Apply date filters if provided
	if dateFrom != "" {
		if fromDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
			query = query.Where("publish_date >= ?", fromDate)
		}
	}
	if dateTo != "" {
		if toDate, err := time.Parse("2006-01-02", dateTo); err == nil {
			query = query.Where("publish_date <= ?", toDate)
		}
	}

	// Apply sorting
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return handleError(ctx, 500, "Error counting records")
	}

	// Get paginated results
	var posts []entity.Post
	if err := query.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return handleError(ctx, 500, "Error fetching posts")
	}

	// Transform to response format
	var postsResponse []response.PostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, response.PostResponse{
			Post_ID:      post.Posts_ID,
			Description:  post.Description,
			Publish_Date: post.Publish_Date,
		})
	}

	// Calculate last page
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return ctx.Status(200).JSON(response.PaginatedPostResponse{
		Data:     postsResponse,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
		PerPage:  limit,
	})
}

// SearchAnnouncementsForProvider - ค้นหาประกาศสำหรับ Provider
func SearchAnnouncementsForProvider(ctx fiber.Ctx) error {
	// ตรวจสอบ role จาก JWT
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)
	username := claims["username"].(string)

	// ตรวจสอบว่าเป็น provider เท่านั้น
	if role != "provider" {
		return handleError(ctx, fiber.StatusUnauthorized, "Unauthorized: requires provider role")
	}

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return handleError(ctx, 404, "Account not found")
	}

	// เพิ่มการตรวจสอบ Provider
	var provider entity.Provider
	if err := database.DB.Where("account_id = ?", account.Account_ID).First(&provider).Error; err != nil {
		return handleError(ctx, 404, "Provider not found")
	}

	// Get search parameters from query
	search := ctx.Query("search")
	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")
	category := ctx.Query("category")
	country := ctx.Query("country")
	educationLevel := ctx.Query("education_level") // เพิ่มการรับค่า education_level
	sortBy := ctx.Query("sort_by", "publish_date")
	sortOrder := ctx.Query("sort_order", "desc")

	// Get pagination parameters
	page, limit, offset := getPaginationParams(ctx)

	// Initialize the base query with provider's account ID filter
	query := database.DB.Model(&entity.Announce_Post{}).
		Where("announce_posts.provider_id = ?", provider.Provider_ID).
		Preload("Category").
		Preload("Country")

	// Apply filters including education_level
	if err := applySearchFilters(query, search, dateFrom, dateTo, category, country, educationLevel, sortBy, sortOrder); err != nil {
		return handleError(ctx, 500, err.Error())
	}

	// Count total records and get results
	total, announcements, err := executeSearchQuery(query, offset, limit)
	if err != nil {
		return handleError(ctx, 500, err.Error())
	}

	// Transform to response format
	announcementsResponse := transformToResponse(announcements)

	// Calculate last page
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return ctx.Status(200).JSON(response.PaginatedAnnouncePostResponse{
		Data:     announcementsResponse,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
		PerPage:  limit,
	})
}

// SearchAnnouncementsForUser - ค้นหาประกาศสำหรับ User (เห็นเฉพาะที่ยังไม่หมดอายุ)
func SearchAnnouncementsForUser(ctx fiber.Ctx) error {
	// Get search parameters
	search := ctx.Query("search")
	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")
	category := ctx.Query("category")
	country := ctx.Query("country")
	educationLevel := ctx.Query("education_level") // เพิ่มการรับค่า education_level
	sortBy := ctx.Query("sort_by", "publish_date")
	sortOrder := ctx.Query("sort_order", "desc")

	// Get pagination parameters
	page, limit, offset := getPaginationParams(ctx)

	// Initialize query for active announcements only
	query := database.DB.Model(&entity.Announce_Post{}).
		Where("(announce_posts.close_date IS NULL OR announce_posts.close_date > ?)", time.Now()).
		Preload("Category").
		Preload("Country")

	// Apply filters including education_level
	if err := applySearchFilters(query, search, dateFrom, dateTo, category, country, educationLevel, sortBy, sortOrder); err != nil {
		return handleError(ctx, 500, err.Error())
	}

	// Count total records and get results
	total, announcements, err := executeSearchQuery(query, offset, limit)
	if err != nil {
		return handleError(ctx, 500, err.Error())
	}

	// Transform to response format
	announcementsResponse := transformToResponse(announcements)

	// Calculate last page
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return ctx.Status(200).JSON(response.PaginatedAnnouncePostResponse{
		Data:     announcementsResponse,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
		PerPage:  limit,
	})
}

// SearchAnnouncementsForAdmin - ค้นหาประกาศสำหรับ Admin และ SuperAdmin (เห็นทั้งหมด)
func SearchAnnouncementsForAdmin(ctx fiber.Ctx) error {
	// ตรวจสอบ role จาก JWT
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)

	// ตรวจสอบว่าเป็น admin หรือ superadmin เท่านั้น
	if role != "admin" && role != "superadmin" {
		return handleError(ctx, fiber.StatusUnauthorized, "Unauthorized: requires admin or superadmin role")
	}

	// Get search parameters
	search := ctx.Query("search")
	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")
	category := ctx.Query("category")
	country := ctx.Query("country")
	educationLevel := ctx.Query("education_level") // เพิ่มการรับค่า education_level
	sortBy := ctx.Query("sort_by", "publish_date")
	sortOrder := ctx.Query("sort_order", "desc")

	// Get pagination parameters
	page, limit, offset := getPaginationParams(ctx)

	// Initialize query for all announcements
	query := database.DB.Model(&entity.Announce_Post{}).
		Preload("Category").
		Preload("Country")

	// Apply filters including education_level
	if err := applySearchFilters(query, search, dateFrom, dateTo, category, country, educationLevel, sortBy, sortOrder); err != nil {
		return handleError(ctx, 500, err.Error())
	}

	// Count total records and get results
	total, announcements, err := executeSearchQuery(query, offset, limit)
	if err != nil {
		return handleError(ctx, 500, err.Error())
	}

	// Transform to response format
	announcementsResponse := transformToResponse(announcements)

	// Calculate last page
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return ctx.Status(200).JSON(response.PaginatedAnnouncePostResponse{
		Data:     announcementsResponse,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
		PerPage:  limit,
	})
}

// Helper functions to reduce code duplication

func applySearchFilters(query *gorm.DB, search, dateFrom, dateTo, category, country, educationLevel, sortBy, sortOrder string) error {
	if search != "" {
		query.Where(
			"LOWER(announce_posts.title) LIKE LOWER(?) OR LOWER(announce_posts.description) LIKE LOWER(?)",
			"%"+search+"%", "%"+search+"%",
		)
	}

	if dateFrom != "" {
		if fromDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
			query.Where("announce_posts.publish_date >= ?", fromDate)
		}
	}

	if dateTo != "" {
		if toDate, err := time.Parse("2006-01-02", dateTo); err == nil {
			query.Where("announce_posts.publish_date <= ?", toDate)
		}
	}

	// Handle multiple categories
	if category != "" {
		categories := strings.Split(category, ",")
		query.Joins("JOIN categories ON announce_posts.category_id = categories.category_id").
			Where("categories.name IN (?)", categories)
	}

	// Handle multiple countries
	if country != "" {
		countries := strings.Split(country, ",")
		query.Joins("JOIN countries ON announce_posts.country_id = countries.country_id").
			Where("countries.name IN (?)", countries)
	}

	// Handle multiple education levels
	if educationLevel != "" {
		levels := strings.Split(educationLevel, ",")
		query.Where("announce_posts.education_level IN (?)", levels)
	}

	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query.Order(fmt.Sprintf("announce_posts.%s %s", sortBy, sortOrder))

	return nil
}

func executeSearchQuery(query *gorm.DB, offset, limit int) (int64, []entity.Announce_Post, error) {
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	var announcements []entity.Announce_Post
	if err := query.Offset(offset).Limit(limit).Find(&announcements).Error; err != nil {
		return 0, nil, err
	}

	return total, announcements, nil
}

func transformToResponse(announcements []entity.Announce_Post) []response.AnnouncePostResponse {
	var announcementsResponse []response.AnnouncePostResponse
	for _, announce := range announcements {
		announcementsResponse = append(announcementsResponse, response.AnnouncePostResponse{
			Announce_ID:     announce.Announce_ID,
			Title:           announce.Title,
			Description:     announce.Description,
			URL:             announce.Url,
			Attach_Name:     announce.Attach_Name,
			Publish_Date:    announce.Publish_Date,
			Close_Date:      announce.Close_Date,
			Category:        announce.Category.Name,
			Country:         announce.Country.Name,
			Education_Level: announce.Education_Level,
			Provider_ID:     announce.Provider_ID, // เปลี่ยนจาก Post_ID เป็น Provider_ID
		})
	}
	return announcementsResponse
}
