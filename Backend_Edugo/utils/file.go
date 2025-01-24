package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
)

const DefaultPathImage = "./public/images/%s"
const DefaultPathAttach = "./public/pdfs/%s"
const MaxFileSize = 5 * 1024 * 1024 // 5MB

func HandleFileImage(ctx fiber.Ctx) error {
	// Handle File Image
	fileImage, errFileImage := ctx.FormFile("image")
	if (errFileImage != nil) {
		log.Println("Error File Image = ", errFileImage)
	}

	var filenameImage *string
	if (fileImage != nil) {
		if (fileImage.Size > MaxFileSize) {
			return ctx.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error message": "File size exceeds 5MB",
			})
		}

		errCheckContentType := checkContentTypeImage(fileImage, "image/jpeg", "image/png", "image/jpg")
		if (errCheckContentType != nil) {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error message": errCheckContentType.Error(),
			})
		}

		// ใช้ filepath.Base เพื่อเก็บเฉพาะชื่อไฟล์
		fileName := filepath.Base(fileImage.Filename)

		// ตรวจสอบว่าชื่อซ้ำหรือไม่
		uniqueFileName := checkUniqueFileName(fmt.Sprintf(DefaultPathImage, fileName))

		// เก็บเฉพาะชื่อไฟล์ใหม่
		filenameOnly := filepath.Base(uniqueFileName)
		filenameImage = &filenameOnly

		// บันทึกไฟล์ใน temp directory
		errSaveFileImage := ctx.SaveFile(fileImage, fmt.Sprintf("./temp/images/%s", filenameOnly))
		if (errSaveFileImage != nil) {
			log.Println("Fail to store file into temp/images directory.")
		}

	} else {
		log.Println("No file uploaded")
		filenameImage = nil
	}
	ctx.Locals("filenameImage", filenameImage)
	return ctx.Next()
}

func HandleFileAttach(ctx fiber.Ctx) error {
	// Handle File Attach
	fileAttach, errFileAttach := ctx.FormFile("attach_file")
	if (errFileAttach != nil) {
		log.Println("Error File Attach = ", errFileAttach)
	}

	var filenameAttach *string
	if (fileAttach != nil) {
		if (fileAttach.Size > MaxFileSize) {
			return ctx.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error message": "File size exceeds 5MB",
			})
		}

		// ตรวจสอบประเภทไฟล์
		errCheckContentType := checkContentTypeImage(fileAttach, "application/pdf")
		if (errCheckContentType != nil) {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error message": errCheckContentType.Error(),
			})
		}

		// ใช้ filepath.Base เพื่อเก็บเฉพาะชื่อไฟล์
		fileName := filepath.Base(fileAttach.Filename)

		// ตรวจสอบชื่อไฟล์ซ้ำ
		uniqueFileName := checkUniqueFileName(fmt.Sprintf(DefaultPathAttach, fileName))

		// เก็บเฉพาะชื่อไฟล์ใหม่
		filenameOnly := filepath.Base(uniqueFileName)
		filenameAttach = &filenameOnly

		// บันทึกไฟล์ใน temp directory
		errSaveFileAttach := ctx.SaveFile(fileAttach, fmt.Sprintf("./temp/pdfs/%s", filenameOnly))
		if (errSaveFileAttach != nil) {
			log.Println("Fail to store file into temp/pdfs directory.")
		}
	} else {
		log.Println("No file uploaded")
		filenameAttach = nil
	}

	// ส่งชื่อไฟล์ไปที่ Locals
	ctx.Locals("filenameAttach", filenameAttach)
	return ctx.Next()
}

func HandleRemoveFileImage(filename string) error {
	err := os.Remove(fmt.Sprintf(DefaultPathImage, filename))
	if (err != nil) {
		log.Println("Failed to remove image file")
		return err
	}
	return nil
}

func HandleRemoveFileAttach(filename string) error {
	err := os.Remove(fmt.Sprintf(DefaultPathAttach, filename))
	if (err != nil) {
		log.Println("Failed to remove attach file")
		return err
	}
	return nil
}

func checkContentTypeImage(file *multipart.FileHeader, contentTypes ...string) error {
	if (len(contentTypes) > 0) {
		for _, contentType := range contentTypes {
			contentTypeFile := file.Header.Get("Content-Type")
			if (contentTypeFile == contentType) {
				return nil
			}
		}
		log.Println("Content-Type File = ", file.Header.Get("Content-Type"))
		return errors.New("not allowed file type")
	} else {
		return errors.New("not found content type")
	}
}

// ClearTempFiles ลบไฟล์ทั้งหมดในโฟลเดอร์ ./temp โดยไม่ลบโฟลเดอร์ ./temp
func ClearTempFiles() {
	// อ่านไฟล์ทั้งหมดในโฟลเดอร์ ./temp
	files, err := os.ReadDir("./temp")
	if (err != nil) {
		log.Println("ไม่สามารถอ่านโฟลเดอร์ ./temp:", err)
		return
	}

	// ลบไฟล์ที่อยู่ในโฟลเดอร์ ./temp
	for _, file := range files {
		filePath := fmt.Sprintf("./temp/%s", file.Name())
		err := os.RemoveAll(filePath) // ลบไฟล์หรือโฟลเดอร์ย่อย
		if (err != nil) {
			log.Printf("ไม่สามารถลบไฟล์หรือโฟลเดอร์ %s: %v", filePath, err)
		} else {
			log.Printf("ลบไฟล์หรือโฟลเดอร์ %s สำเร็จ", filePath)
		}
	}
}

func CreateTempFolder() {
	err := os.MkdirAll("./temp/images", 0755)
	if (err != nil) {
		log.Println("Failed to create temp/images directory:", err)
	}

	err = os.MkdirAll("./temp/pdfs", 0755)
	if (err != nil) {
		log.Println("Failed to create temp/pdfs directory:", err)
	}
}

// RemoveTempToPublic ย้ายไฟล์จากโฟลเดอร์ temp ไปยังโฟลเดอร์ public
func RemoveTempToPublic() {
	// ย้ายไฟล์รูปภาพจาก ./temp/images ไปที่ ./public/images
	moveFiles("./temp/images", "./public/images")

	// ย้ายไฟล์ PDF จาก ./temp/pdfs ไปที่ ./public/pdfs
	moveFiles("./temp/pdfs", "./public/pdfs")
}

// moveFiles ย้ายไฟล์ทั้งหมดจาก sourceDir ไปที่ destDir
func moveFiles(sourceDir, destDir string) {
	files, err := os.ReadDir(sourceDir)
	if (err != nil) {
		log.Printf("ไม่สามารถอ่านโฟลเดอร์ %s: %v", sourceDir, err)
		return
	}

	for _, file := range files {
		sourcePath := fmt.Sprintf("%s/%s", sourceDir, file.Name())
		destPath := fmt.Sprintf("%s/%s", destDir, file.Name())

		// ย้ายไฟล์จาก sourcePath ไปที่ destPath
		err := os.Rename(sourcePath, destPath)
		if (err != nil) {
			log.Printf("ไม่สามารถย้ายไฟล์ %s ไปที่ %s: %v", sourcePath, destPath, err)
		} else {
			log.Printf("ย้ายไฟล์ %s ไปที่ %s สำเร็จ", sourcePath, destPath)
		}
	}
}

func checkUniqueFileName(filePath string) string {
	// แยกชื่อไฟล์และนามสกุล
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	// ตรวจสอบว่าไฟล์มีอยู่หรือไม่ ถ้ามีให้เพิ่มตัวเลขต่อท้าย
	counter := 1
	newFilePath := filePath
	for fileExists(newFilePath) {
		newFilePath = fmt.Sprintf("%s/%s-%d%s", dir, name, counter, ext)
		counter++
	}
	return newFilePath
}

// ฟังก์ชันตรวจสอบว่าไฟล์มีอยู่หรือไม่
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
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
