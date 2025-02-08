package utils

import (
	"errors"
	"io"
	"log"
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
)

const MaxFileSize = 5 * 1024 * 1024 // 5MB

func checkContentTypeImage(file *multipart.FileHeader, contentTypes ...string) error {
	if len(contentTypes) > 0 {
		for _, contentType := range contentTypes {
			contentTypeFile := file.Header.Get("Content-Type")
			if contentTypeFile == contentType {
				return nil
			}
		}
		log.Println("Content-Type File = ", file.Header.Get("Content-Type"))
		return errors.New("not allowed file type")
	} else {
		return errors.New("not found content type")
	}
}

// เพิ่มฟังก์ชันใหม่ที่รับพารามิเตอร์ fieldName
func HandleImageUpload(ctx fiber.Ctx, fieldName string) error {
	// Handle panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in HandleImageUpload: %v", r)
			ctx.Locals("imageBytes", nil)
		}
	}()

	// Try to get the file from form
	file, err := ctx.FormFile(fieldName)

	// If no file is uploaded, set imageBytes to nil and return nil
	if err != nil {
		ctx.Locals("imageBytes", nil)
		return nil // เปลี่ยนจาก ctx.Next() เป็น return nil
	}

	// If file exists, process it
	if file != nil {
		// ตรวจสอบขนาดไฟล์
		if file.Size > MaxFileSize {
			return errors.New("file size exceeds 5MB")
		}

		// ตรวจสอบประเภทไฟล์
		if err := checkContentTypeImage(file, "image/jpeg", "image/png", "image/jpg"); err != nil {
			return err
		}

		fileContent, err := file.Open()
		if err != nil {
			ctx.Locals("imageBytes", nil)
			return err
		}
		defer fileContent.Close()

		imageBytes, err := io.ReadAll(fileContent)
		if err != nil {
			ctx.Locals("imageBytes", nil)
			return err
		}

		ctx.Locals("imageBytes", imageBytes)
	} else {
		ctx.Locals("imageBytes", nil)
	}

	return nil // เปลี่ยนจาก ctx.Next() เป็น return nil
}

// เพิ่มฟังก์ชันใหม่ที่รับพารามิเตอร์ fieldName
func HandleAttachUpload(ctx fiber.Ctx, fieldName string) error {
	// Handle panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in HandleAttachUpload: %v", r)
			ctx.Locals("attachBytes", nil)
		}
	}()

	// Try to get the file from form
	file, err := ctx.FormFile(fieldName)

	// If no file is uploaded, set attachBytes to nil and return nil
	if err != nil {
		ctx.Locals("attachBytes", nil)
		return nil // เปลี่ยนจาก ctx.Next() เป็น return nil
	}

	// If file exists, process it
	if file != nil {
		// ตรวจสอบขนาดไฟล์
		if file.Size > MaxFileSize {
			return errors.New("file size exceeds 5MB")
		}

		// ตรวจสอบประเภทไฟล์
		if err := checkContentTypeImage(file, "application/pdf"); err != nil {
			return err
		}

		fileContent, err := file.Open()
		if err != nil {
			ctx.Locals("attachBytes", nil)
			return err
		}
		defer fileContent.Close()

		attachBytes, err := io.ReadAll(fileContent)
		if err != nil {
			ctx.Locals("attachBytes", nil)
			return err
		}

		ctx.Locals("attachBytes", attachBytes)
	} else {
		ctx.Locals("attachBytes", nil)
	}

	return nil // เปลี่ยนจาก ctx.Next() เป็น return nil
}

// เพิ่มฟังก์ชันใหม่ที่รับพารามิเตอร์ fieldName
func HandleAvatarUpload(ctx fiber.Ctx, fieldName string) error {
	// Handle panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in HandleAvatarUpload: %v", r)
			ctx.Locals("avatarBytes", nil)
		}
	}()

	// Try to get the file from form
	file, err := ctx.FormFile(fieldName)

	// If no file is uploaded, set avatarBytes to nil and return nil
	if err != nil {
		ctx.Locals("avatarBytes", nil)
		return nil
	}

	// If file exists, process it
	if file != nil {
		// Check file size (limit to 2MB for avatars)
		if file.Size > (2 * 1024 * 1024) {
			return errors.New("avatar file size exceeds 2MB")
		}

		// Check file type (only allow common image formats)
		if err := checkContentTypeImage(file, "image/jpeg", "image/png", "image/jpg"); err != nil {
			return err
		}

		fileContent, err := file.Open()
		if err != nil {
			ctx.Locals("avatarBytes", nil)
			return err
		}
		defer fileContent.Close()

		avatarBytes, err := io.ReadAll(fileContent)
		if err != nil {
			ctx.Locals("avatarBytes", nil)
			return err
		}

		ctx.Locals("avatarBytes", avatarBytes)
	} else {
		ctx.Locals("avatarBytes", nil)
	}

	return nil
}
