package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/utils"
	"strings" // added import
)

// AuthProvider ใช้ตรวจสอบการยืนยันตัวตนของผู้ใช้ผ่าน JWT
func AuthProvider(ctx fiber.Ctx) error {
	// ดึงค่า Authorization จาก headers
	token := ctx.Get("Authorization")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	// แยกโทเค็นออกเป็นส่วนๆ
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	token = parts[1]
	// ถอดรหัสโทเค็นเพื่อนำข้อมูลผู้ใช้มาใช้
	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	
	// ตรวจสอบบทบาทของผู้ใช้
	role := claims["role"].(string)
	if role != "provider" && role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden access",
		})
	}

	// ดำเนินการต่อกับ middleware ถัดไป
	// ctx.Locals("user", claims)
	// ctx.Locals("role", claims["role"])
	return ctx.Next()
}

// PermissionCreate เป็น middleware ที่อนุญาตให้คำขอดำเนินต่อไปทันที
func PermissionCreate(ctx fiber.Ctx) error {
	return ctx.Next()
}