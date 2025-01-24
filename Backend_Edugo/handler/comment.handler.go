package handler

import (
	"fmt"
	"io"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func CreateComment(ctx fiber.Ctx) error {
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

	// check available account
	var account entity.Account
	result = database.DB.Where("account_id = ?", comment.Account_ID).First(&account)
	if result.RowsAffected == 0 {
		return utils.HandleError(ctx, 400, "Account not found")
	}

	// Create comment
	newComment := entity.Comment{
		Comments_Text: comment.Comments_Text,
		Posts_ID:      comment.Posts_ID,
		Account_ID:    comment.Account_ID,
	}

	// จัดการไฟล์รูปภาพ
	file, err := ctx.FormFile("comments_image")
	if err == nil && file != nil {
		// เปิดไฟล์
		fileContent, err := file.Open()
		if err != nil {
			return utils.HandleError(ctx, 400, "Cannot open image file: "+err.Error())
		}
		defer fileContent.Close()

		// อ่านข้อมูลไฟล์
		imageBytes, err := io.ReadAll(fileContent)
		if err != nil {
			return utils.HandleError(ctx, 400, "Error reading image file: "+err.Error())
		}

		newComment.Comments_Image = imageBytes
		fmt.Printf("Image size: %d bytes\n", len(imageBytes))
	}

	// Create comment with debug log
	if err := database.DB.Debug().Create(&newComment).Error; err != nil {
		return utils.HandleError(ctx, 400, "Error saving comment: "+err.Error())
	}

	// Return response
	commentResponse := response.CommentResponse{
		Comments_ID:     newComment.Comments_ID,
		Comments_Text:  newComment.Comments_Text,
		Publish_Date:   newComment.Publish_Date,
		Posts_ID:       newComment.Posts_ID,
		Account_ID:     newComment.Account_ID,
	}

	return ctx.Status(201).JSON(commentResponse)
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
