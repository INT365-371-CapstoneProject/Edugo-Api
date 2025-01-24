package utils

import (
	"github.com/gofiber/fiber/v3"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
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

// HandleError - ฟังก์ชันสำหรับจัดการข้อผิดพลาด (exported function)
func HandleError(ctx fiber.Ctx, statusCode int, message string) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
}
