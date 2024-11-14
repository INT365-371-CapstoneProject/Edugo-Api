package utils

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
)

var im int = 1
var pdf int = 1

const DefaultPathImage = "./public/images/%s"
const DefaultPathAttach = "./public/pdfs/%s"

func HandleFileImage(ctx fiber.Ctx) error {
	// Handle File Image
	fileImage, errFileImage := ctx.FormFile("image")
	if errFileImage != nil {
		log.Println("Error File Image = ", errFileImage)
	}

	var filenameImage *string
	if fileImage != nil {
		errCheckContentType := checkContentTypeImage(fileImage, "image/jpg", "image/png")
		if errCheckContentType != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error message": errCheckContentType.Error(),
			})
		}

		filenameImage = &fileImage.Filename
		newFilenameImage := fmt.Sprintf("%d", im)
		errSaveFileImage := ctx.SaveFile(fileImage, fmt.Sprintf("./temp/images/%s", newFilenameImage))
		if errSaveFileImage != nil {
			log.Println("Fail to store file into temp/images directory.")
		} else {
			im++
			filenameImage = &newFilenameImage
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
	if errFileAttach != nil {
		log.Println("Error File Attach = ", errFileAttach)
	}

	var filenameAttach *string
	if fileAttach != nil {
		errCheckContentType := checkContentTypeImage(fileAttach, "application/pdf")
		if errCheckContentType != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error message": errCheckContentType.Error(),
			})
		}

		filenameAttach = &fileAttach.Filename
		newFilenameAttach := fmt.Sprintf("%d", pdf)
		errSaveFileAttach := ctx.SaveFile(fileAttach, fmt.Sprintf("./temp/pdfs/%s", newFilenameAttach))
		if errSaveFileAttach != nil {
			log.Println("Fail to store file into temp/pdfs directory.")
		} else {
			pdf++
			filenameAttach = &newFilenameAttach
		}
	} else {
		log.Println("No file uploaded")
		filenameAttach = nil
	}
	ctx.Locals("filenameAttach", filenameAttach)
	return ctx.Next()
}

func HandleRemoveFileImage(filename string) error {
	err := os.Remove(fmt.Sprintf(DefaultPathImage, filename))
	if err != nil {
		log.Println("Failed to remove image file")
		return err
	}
	return nil
}

func HandleRemoveFileAttach(filename string) error {
	err := os.Remove(fmt.Sprintf(DefaultPathAttach, filename))
	if err != nil {
		log.Println("Failed to remove attach file")
		return err
	}
	return nil
}

func checkContentTypeImage(file *multipart.FileHeader, contentTypes ...string) error {
	if len(contentTypes) > 0 {
		for _, contentType := range contentTypes {
			contentTypeFile := file.Header.Get("Content-Type")
			if contentTypeFile == contentType {
				return nil
			}
		}
		return errors.New("not allowed file type")
	} else {
		return errors.New("not found content type")
	}
}

// ClearTempFiles ลบไฟล์ทั้งหมดในโฟลเดอร์ ./temp โดยไม่ลบโฟลเดอร์ ./temp
func ClearTempFiles() {
	// อ่านไฟล์ทั้งหมดในโฟลเดอร์ ./temp
	files, err := os.ReadDir("./temp")
	if err != nil {
		log.Println("ไม่สามารถอ่านโฟลเดอร์ ./temp:", err)
		return
	}

	// ลบไฟล์ที่อยู่ในโฟลเดอร์ ./temp
	for _, file := range files {
		filePath := fmt.Sprintf("./temp/%s", file.Name())
		err := os.RemoveAll(filePath) // ลบไฟล์หรือโฟลเดอร์ย่อย
		if err != nil {
			log.Printf("ไม่สามารถลบไฟล์หรือโฟลเดอร์ %s: %v", filePath, err)
		} else {
			log.Printf("ลบไฟล์หรือโฟลเดอร์ %s สำเร็จ", filePath)
		}
	}
}

func CreateTempFolder() {
	err := os.MkdirAll("./temp/images", 0755)
	if err != nil {
		log.Println("Failed to create temp/images directory:", err)
	}

	err = os.MkdirAll("./temp/pdfs", 0755)
	if err != nil {
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
	if err != nil {
		log.Printf("ไม่สามารถอ่านโฟลเดอร์ %s: %v", sourceDir, err)
		return
	}

	for _, file := range files {
		sourcePath := fmt.Sprintf("%s/%s", sourceDir, file.Name())
		destPath := fmt.Sprintf("%s/%s", destDir, file.Name())

		// ตรวจสอบว่าไฟล์ใน destPath มีชื่อเดียวกันหรือไม่
		if fileExists(destPath) {
			// หากไฟล์มีชื่อเดียวกัน, เปลี่ยนชื่อไฟล์ใหม่โดยเพิ่มตัวเลข
			destPath = generateNewFileName(destPath)
		}

		// ย้ายไฟล์จาก sourcePath ไปที่ destPath
		err := os.Rename(sourcePath, destPath)
		if err != nil {
			log.Printf("ไม่สามารถย้ายไฟล์ %s ไปที่ %s: %v", sourcePath, destPath, err)
		} else {
			log.Printf("ย้ายไฟล์ %s ไปที่ %s สำเร็จ", sourcePath, destPath)
		}
	}
}

// ตรวจสอบว่าไฟล์ใน destPath มีอยู่หรือไม่
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// สร้างชื่อไฟล์ใหม่โดยเพิ่มตัวเลขที่ปลายชื่อไฟล์
func generateNewFileName(destPath string) string {
	ext := fmt.Sprintf(".%s", getFileExtension(destPath))  // เอานามสกุลไฟล์
	baseName := strings.TrimSuffix(destPath, ext)          // ชื่อไฟล์หลักโดยไม่รวมส่วนขยาย

	// เริ่มจากเลข 1 และเพิ่มขึ้นเรื่อยๆ
	counter := 1
	newDestPath := fmt.Sprintf("%s-%d%s", baseName, counter, ext)

	// ตรวจสอบว่าชื่อไฟล์ใหม่ซ้ำหรือไม่ ถ้าซ้ำจะเพิ่มตัวเลขไปเรื่อยๆ
	for fileExists(newDestPath) {
		counter++
		newDestPath = fmt.Sprintf("%s-%d%s", baseName, counter, ext)
	}

	return newDestPath
}

// ฟังก์ชันนี้ใช้เพื่อดึงนามสกุลของไฟล์
func getFileExtension(fileName string) string {
	ext := strings.ToLower(fileName[strings.LastIndex(fileName, ".")+1:])
	return ext
}
