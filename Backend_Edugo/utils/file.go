package utils

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"

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
		errSaveFileImage := ctx.SaveFile(fileImage, fmt.Sprintf("./public/images/%s", newFilenameImage))
		if errSaveFileImage != nil {
			log.Println("Fail to store file into public/images directory.")
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
