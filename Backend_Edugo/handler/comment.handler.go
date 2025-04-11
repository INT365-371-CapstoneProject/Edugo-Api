package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/middleware"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func CreateComment(ctx fiber.Ctx) error {

	claims := middleware.GetTokenClaims(ctx)
	username := claims["username"].(string)

	// หา account จาก username
	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return handleError(ctx, 404, "Account not found")
	}

	comment := new(request.CreateCommentRequest)
	if err := ctx.Bind().Body(comment); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(comment); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// check available post
	var post entity.Post
	result := database.DB.Where("posts_id = ?", comment.Posts_ID).First(&post)
	if result.RowsAffected == 0 {
		return utils.HandleError(ctx, 400, "Post not found")
	}

	// ใช้ฟังก์ชัน HandleImageUpload แทนการจัดการไฟล์โดยตรง
	if err := utils.HandleImageUpload(ctx, "comments_image"); err != nil {
		return utils.HandleError(ctx, 400, "Error handling image upload: "+err.Error())
	}

	// Create comment
	newComment := entity.Comment{
		Comments_Text: comment.Comments_Text,
		Posts_ID:      comment.Posts_ID,
		Account_ID:    account.Account_ID,
	}

	// ตรวจสอบว่ามีรูปภาพถูกอัพโหลดหรือไม่
	if imageBytes := ctx.Locals("imageBytes"); imageBytes != nil {
		newComment.Comments_Image = imageBytes.([]byte)
	}

	// Create comment with debug log
	if err := database.DB.Debug().Create(&newComment).Error; err != nil {
		return utils.HandleError(ctx, 400, "Error saving comment: "+err.Error())
	}

	// Return response
	commentResponse := response.CommentResponse{
		Comments_ID:   newComment.Comments_ID,
		Comments_Text: newComment.Comments_Text,
		Publish_Date:  newComment.Publish_Date,
		Posts_ID:      newComment.Posts_ID,
		Account_ID:    newComment.Account_ID,
	}

	return ctx.Status(201).JSON(commentResponse)
}

func GetAllComment(ctx fiber.Ctx) error {
	var comments []struct {
		entity.Comment
		Fullname string `json:"fullname"`
	}
	result := database.DB.Table("comments c").
		Select(`c.comments_id, c.comments_text, c.publish_date, c.posts_id, c.account_id,
			CASE
				WHEN pr.company_name IS NOT NULL THEN pr.company_name
				ELSE CONCAT(a.first_name, ' ', a.last_name)
			END AS fullname`).
		Joins("JOIN accounts a ON c.account_id = a.account_id").
		Joins("LEFT JOIN providers pr ON a.account_id = pr.account_id").
		Scan(&comments)

	if result.Error != nil {
		return utils.HandleError(ctx, 500, "Error retrieving comments: "+result.Error.Error())
	}

	// สร้าง response list
	var commentResponses []response.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, response.CommentResponse{
			Comments_ID:   comment.Comments_ID,
			Comments_Text: comment.Comments_Text,
			Publish_Date:  comment.Publish_Date,
			Posts_ID:      comment.Posts_ID,
			Account_ID:    comment.Account_ID,
			Fullname:      comment.Fullname,
		})
	}

	return ctx.Status(200).JSON(commentResponses)
}

func GetCommentByPostID(ctx fiber.Ctx) error {
	postID := ctx.Params("post_id")
	var comments []struct {
		entity.Comment
		Fullname string `json:"fullname"`
	}

	result := database.DB.Table("comments c").
		Select(`c.comments_id, c.comments_text, c.publish_date, c.posts_id, c.account_id,
			CASE
				WHEN pr.company_name IS NOT NULL THEN pr.company_name
				ELSE CONCAT(a.first_name, ' ', a.last_name)
			END AS fullname`).
		Joins("JOIN accounts a ON c.account_id = a.account_id").
		Joins("LEFT JOIN providers pr ON a.account_id = pr.account_id").
		Where("c.posts_id = ?", postID).
		Scan(&comments)

	if len(comments) == 0 {
		return utils.HandleError(ctx, 404, "No comments found for this post")
	}

	if result.Error != nil {
		return utils.HandleError(ctx, 500, "Error retrieving comments: "+result.Error.Error())
	}

	// สร้าง response list
	var commentResponses []response.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, response.CommentResponse{
			Comments_ID:   comment.Comments_ID,
			Comments_Text: comment.Comments_Text,
			Publish_Date:  comment.Publish_Date,
			Posts_ID:      comment.Posts_ID,
			Account_ID:    comment.Account_ID,
			Fullname:      comment.Fullname,
		})
	}

	return ctx.Status(200).JSON(commentResponses)
}

func GetCommentAvatarImageByAccountID(ctx fiber.Ctx) error {
	commentID := ctx.Params("id")
	var comment entity.Comment

	if err := database.DB.Where("comments_id = ?", commentID).First(&comment).Error; err != nil {
		return utils.HandleError(ctx, 404, "Comment not found")
	}

	var account entity.Account
	if err := database.DB.Select("avatar").First(&account, "account_id = ?", comment.Account_ID).Error; err != nil {
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

func GetCommentImage(ctx fiber.Ctx) error {
	commentID := ctx.Params("id")
	var comment entity.Comment

	if err := database.DB.Where("comments_id = ?", commentID).First(&comment).Error; err != nil {
		return utils.HandleError(ctx, 404, "Comment not found")
	}

	if comment.Comments_Image == nil {
		return utils.HandleError(ctx, 404, "No image found for this comment")
	}

	ctx.Set("Content-Type", "image/jpeg") // เปลี่ยนเป็นประเภทภาพที่ถูกต้อง
	return ctx.Send(comment.Comments_Image)
}
