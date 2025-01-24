package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/utils"
	"strings" // added import
)

// extractAndDecodeToken extracts and decodes the JWT token from the Authorization header
func extractAndDecodeToken(ctx fiber.Ctx) (map[string]interface{}, error) {
	token := ctx.Get("Authorization")
	if token == "" {
		return nil, errors.New("Unauthorized")
	}
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("Unauthorized")
	}
	return utils.DecodeToken(parts[1])
}

// AuthProvider ใช้ตรวจสอบการยืนยันตัวตนของผู้ใช้ผ่าน JWT
func AuthProvider(ctx fiber.Ctx) error {
	claims, err := extractAndDecodeToken(ctx)
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
	ctx.Locals("user", claims)
	return ctx.Next()
}

// AuthAny ใช้ตรวจสอบการยืนยันตัวตนของผู้ใช้ผ่าน JWT โดยไม่ตรวจสอบบทบาท
func AuthAny(ctx fiber.Ctx) error {
	claims, err := extractAndDecodeToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// ดำเนินการต่อกับ middleware ถัดไป
	ctx.Locals("user", claims)
	return ctx.Next()
}

// PermissionCreate เป็น middleware ที่อนุญาตให้คำขอดำเนินต่อไปทันที
func PermissionCreate(ctx fiber.Ctx) error {
	return ctx.Next()
}

func GetTokenClaims(ctx fiber.Ctx) map[string]interface{} {
	claims, _ := extractAndDecodeToken(ctx)
	return claims
}