package handler

import (
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/middleware"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
)

func CreateAdminForSuperadmin(ctx fiber.Ctx) error {
	admin := new(request.AdminCreateRequest)
	if err := ctx.Bind().Body(admin); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(admin); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// check duplicate email
	var account entity.Account
	result := database.DB.Where("email = ?", admin.Email).First(&account)
	if result.RowsAffected > 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// check duplicate username
	result = database.DB.Where("username = ?", admin.Username).First(&account)
	if result.RowsAffected > 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Username already exists",
		})
	}

	// Begin transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return ctx.Status(409).JSON(fiber.Map{
			"message": "Failed to begin transaction",
		})
	}

	// Create account
	newAccount := entity.Account{
		Username:   admin.Username,
		Email:      admin.Email,
		FirstName:  &admin.FirstName,
		LastName:   &admin.LastName,
		Status:     "Active",
		Last_Login: nil,
		Role:       "admin",
	}

	// Hash password
	hashedPassword, err := utils.HashingPassword(admin.Password)
	if err != nil {
		tx.Rollback()
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}
	newAccount.Password = hashedPassword

	// Insert account to database
	if err := tx.Create(&newAccount).Error; err != nil {
		tx.Rollback()
		return ctx.Status(409).JSON(fiber.Map{
			"message": "Failed to create account",
		})
	}

	// Create admin - แก้ไขส่วนนี้
	newAdmin := entity.Admin{
		Account_ID: newAccount.Account_ID, // ใช้ Account_ID จาก account ที่เพิ่งสร้าง
		Phone:      admin.Phone,
	}

	// Insert admin to database
	if err := tx.Create(&newAdmin).Error; err != nil {
		tx.Rollback()
		return ctx.Status(409).JSON(fiber.Map{
			"message": "Failed to create admin",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return ctx.Status(409).JSON(fiber.Map{
			"message": "Failed to commit transaction",
		})
	}

	// Create response
	providerResponse := response.AdminResponse{
		Admin_ID:  newAccount.Account_ID,
		Username:  newAccount.Username,
		Email:     newAccount.Email,
		FirstName: *newAccount.FirstName,
		LastName:  *newAccount.LastName,
		Status:    newAccount.Status,
		Phone:     newAdmin.Phone,
	}

	return ctx.Status(201).JSON(providerResponse)
}

func GetAllProviderForAdmin(ctx fiber.Ctx) error {
	// รับค่า query parameter สำหรับ pagination
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	// คำนวณ offset
	offset := (page - 1) * limit

	// นับจำนวน provider ทั้งหมด
	var totalCount int64
	database.DB.Model(&entity.Provider{}).Count(&totalCount)

	// คำนวณจำนวนหน้าทั้งหมด
	totalPages := math.Ceil(float64(totalCount) / float64(limit))

	// ดึงข้อมูล provider ตาม pagination
	var providers []entity.Provider
	result := database.DB.Preload("Account", "role = ?", "provider").
		Offset(offset).
		Limit(limit).
		Find(&providers)

	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}

	var providerResponse []response.ProviderResponse
	for _, provider := range providers {
		providerResponse = append(providerResponse, response.ProviderResponse{
			Provider_ID:  provider.Provider_ID,
			Company_Name: provider.Company_Name,
			FirstName:    provider.Account.FirstName,
			LastName:     provider.Account.LastName,
			Username:     provider.Account.Username,
			Email:        provider.Account.Email,
			URL:          provider.URL,
			Address:      provider.Address,
			City:         provider.City,
			Country:      provider.Country,
			Postal_Code:  provider.Postal_Code,
			Phone:        provider.Phone,
			Status:       provider.Account.Status,
			Verify:       provider.Verify,
			Create_On:    provider.Account.Create_On,
			Last_Login:   provider.Account.Last_Login,
			Update_On:    provider.Account.Update_On,
			Role:         provider.Account.Role,
		})
	}

	// ส่งข้อมูล pagination กลับไปพร้อมกับข้อมูล provider
	return ctx.Status(200).JSON(fiber.Map{
		"providers": providerResponse,
		"pagination": fiber.Map{
			"total":      totalCount,
			"page":       page,
			"limit":      limit,
			"total_page": totalPages,
		},
	})
}

func GetIDProviderForAdmin(ctx fiber.Ctx) error {
	// Get provider ID from params
	providerID := ctx.Params("id")

	// Find provider in database with Account preloaded
	var provider entity.Provider
	result := database.DB.Preload("Account", "role = ?", "provider").First(&provider, providerID)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Provider not found",
		})
	}

	// Create response
	providerResponse := response.ProviderResponse{
		Provider_ID:  provider.Provider_ID,
		Company_Name: provider.Company_Name,
		FirstName:    provider.Account.FirstName,
		LastName:     provider.Account.LastName,
		Username:     provider.Account.Username,
		Email:        provider.Account.Email,
		URL:          provider.URL,
		Address:      provider.Address,
		City:         provider.City,
		Country:      provider.Country,
		Postal_Code:  provider.Postal_Code,
		Phone:        provider.Phone,
		Status:       provider.Account.Status,
		Verify:       provider.Verify,
		Create_On:    provider.Account.Create_On,
		Last_Login:   provider.Account.Last_Login,
		Update_On:    provider.Account.Update_On,
		Role:         provider.Account.Role,
	}

	return ctx.Status(200).JSON(providerResponse)
}

func VerifyProviderForAdmin(ctx fiber.Ctx) error {
	// Get provider ID from params
	providerID := ctx.Params("id")

	// Bind request body to get verification status
	verifyRequest := struct {
		Status string `json:"status" validate:"required,oneof=Yes No"`
	}{}

	if err := ctx.Bind().Body(&verifyRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request body: " + err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(verifyRequest); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// Find provider in database
	var provider entity.Provider
	result := database.DB.First(&provider, providerID)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Provider not found",
		})
	}

	// ใช้ค่าที่รับมาจาก frontend โดยตรง เนื่องจากค่า ENUM ในฐานข้อมูลตรงกับค่าที่รับมา
	// Update provider verify status
	if err := database.DB.Model(&provider).Update("verify", verifyRequest.Status).Error; err != nil {
		return ctx.Status(409).JSON(fiber.Map{
			"message": "Failed to update provider verification status",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Provider verification status updated to " + verifyRequest.Status,
	})
}

func GetAllUser(ctx fiber.Ctx) error {
	// ตรวจสอบ role จาก JWT
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)
	username := claims["username"].(string)

	// รับค่า query parameter สำหรับ pagination
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	// คำนวณ offset
	offset := (page - 1) * limit

	// เตรียมตัวแปรสำหรับเก็บข้อมูลผู้ใช้
	var users []entity.Account
	var totalCount int64
	var userResponse []response.UserResponse

	// เงื่อนไขการดึงข้อมูลตาม role
	if role == "superadmin" {
		// ถ้าเป็น superadmin ดึงข้อมูล user, provider และ admin
		// นับจำนวนข้อมูลทั้งหมด
		database.DB.Model(&entity.Account{}).Where("role IN ?", []string{"user", "provider", "admin"}).Count(&totalCount)

		// ดึงข้อมูลตาม pagination และเรียงตาม role
		if err := database.DB.Where("role IN ?", []string{"user", "provider", "admin"}).
			Order("CASE " +
				"WHEN role = 'admin' THEN 1 " +
				"WHEN role = 'provider' THEN 2 " +
				"WHEN role = 'user' THEN 3 " +
				"ELSE 4 END").
			Offset(offset).
			Limit(limit).
			Find(&users).Error; err != nil {
			return ctx.Status(409).JSON(fiber.Map{
				"message": "Failed to fetch users: " + err.Error(),
			})
		}
	} else if role == "admin" {
		// ถ้าเป็น admin ดึงข้อมูลเฉพาะ user และ provider
		// นับจำนวนข้อมูลทั้งหมด
		database.DB.Model(&entity.Account{}).Where("role IN ?", []string{"user", "provider"}).Count(&totalCount)

		// ดึงข้อมูลตาม pagination และเรียงตาม role
		if err := database.DB.Where("role IN ?", []string{"user", "provider"}).
			Order("CASE " +
				"WHEN role = 'provider' THEN 1 " +
				"WHEN role = 'user' THEN 2 " +
				"ELSE 3 END").
			Offset(offset).
			Limit(limit).
			Find(&users).Error; err != nil {
			return ctx.Status(409).JSON(fiber.Map{
				"message": "Failed to fetch users: " + err.Error(),
			})
		}
	} else {
		// ถ้าไม่ใช่ admin หรือ superadmin จะไม่มีสิทธิ์เข้าถึง
		return ctx.Status(403).JSON(fiber.Map{
			"message": "Access denied: Insufficient permissions",
		})
	}

	// คำนวณจำนวนหน้าทั้งหมด
	totalPages := math.Ceil(float64(totalCount) / float64(limit))

	// แปลงข้อมูลเป็น response format
	for _, user := range users {
		var firstName, lastName string
		if user.FirstName != nil {
			firstName = *user.FirstName
		}
		if user.LastName != nil {
			lastName = *user.LastName
		}

		// สร้าง response object แต่ละ user
		userResp := response.UserResponse{
			Account_ID: user.Account_ID,
			Username:   user.Username,
			Email:      user.Email,
			FirstName:  &firstName,
			LastName:   &lastName,
			Status:     user.Status,
			Role:       user.Role,
			Create_On:  user.Create_On,
			Last_Login: user.Last_Login,
			Update_On:  user.Update_On,
		}

		// เพิ่มข้อมูลเพิ่มเติมตาม role
		if user.Role == "provider" {
			var provider entity.Provider
			if err := database.DB.Where("account_id = ?", user.Account_ID).First(&provider).Error; err == nil {
				// ถ้าเป็น provider เพิ่มข้อมูลเฉพาะของ provider
				userResp.ProviderDetails = &response.ProviderDetails{
					Provider_ID:  provider.Provider_ID,
					Company_Name: provider.Company_Name,
					URL:          provider.URL,
					Address:      provider.Address,
					City:         provider.City,
					Country:      provider.Country,
					Postal_Code:  provider.Postal_Code,
					Phone:        provider.Phone,
					Phone_Person: provider.Phone_Person,
					Verify:       provider.Verify,
				}
			}
		} else if user.Role == "admin" {
			var admin entity.Admin
			if err := database.DB.Where("account_id = ?", user.Account_ID).First(&admin).Error; err == nil {
				// ถ้าเป็น admin เพิ่มข้อมูลเฉพาะของ admin
				userResp.AdminDetails = &response.AdminDetails{
					Admin_ID: admin.Admin_ID,
					Phone:    admin.Phone,
				}
			}
		}

		userResponse = append(userResponse, userResp)
	}

	// ส่งข้อมูลผู้ใช้ทั้งหมดกลับไปพร้อมข้อมูล pagination
	return ctx.Status(200).JSON(fiber.Map{
		"users":    userResponse,
		"count":    len(userResponse),
		"role":     role,
		"username": username,
		"pagination": fiber.Map{
			"total":      totalCount,
			"page":       page,
			"limit":      limit,
			"total_page": totalPages,
		},
	})
}

// GetIDUser - ฟังก์ชันสำหรับดึงข้อมูลผู้ใช้จากไอดี
func GetIDUser(ctx fiber.Ctx) error {
	// ตรวจสอบสิทธิ์จาก JWT token
	claims := middleware.GetTokenClaims(ctx)
	role := claims["role"].(string)

	// ตรวจสอบว่ามีสิทธิ์เข้าถึงหรือไม่
	if role != "superadmin" && role != "admin" {
		return ctx.Status(403).JSON(fiber.Map{
			"message": "ไม่มีสิทธิ์เข้าถึงข้อมูลนี้",
		})
	}

	// รับค่า ID จาก parameter
	userID := ctx.Params("id")

	// ค้นหาข้อมูลบัญชีผู้ใช้
	var user entity.Account
	if err := database.DB.First(&user, userID).Error; err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "ไม่พบข้อมูลผู้ใช้",
		})
	}

	// ตรวจสอบสิทธิ์ในการเข้าถึงข้อมูล
	if role == "admin" && (user.Role == "admin" || user.Role == "superadmin") {
		return ctx.Status(403).JSON(fiber.Map{
			"message": "ไม่มีสิทธิ์เข้าถึงข้อมูลของแอดมินหรือซูเปอร์แอดมิน",
		})
	}

	// สร้าง response
	var firstName, lastName string
	if user.FirstName != nil {
		firstName = *user.FirstName
	}
	if user.LastName != nil {
		lastName = *user.LastName
	}

	userResp := response.UserResponse{
		Account_ID: user.Account_ID,
		Username:   user.Username,
		Email:      user.Email,
		FirstName:  &firstName,
		LastName:   &lastName,
		Status:     user.Status,
		Role:       user.Role,
		Create_On:  user.Create_On,
		Last_Login: user.Last_Login,
		Update_On:  user.Update_On,
	}

	// เพิ่มข้อมูลเพิ่มเติมตาม role
	if user.Role == "provider" {
		var provider entity.Provider
		if err := database.DB.Where("account_id = ?", user.Account_ID).First(&provider).Error; err == nil {
			// ถ้าเป็น provider เพิ่มข้อมูลเฉพาะของ provider
			userResp.ProviderDetails = &response.ProviderDetails{
				Provider_ID:  provider.Provider_ID,
				Company_Name: provider.Company_Name,
				URL:          provider.URL,
				Address:      provider.Address,
				City:         provider.City,
				Country:      provider.Country,
				Postal_Code:  provider.Postal_Code,
				Phone:        provider.Phone,
				Phone_Person: provider.Phone_Person,
				Verify:       provider.Verify,
			}
		}
	} else if user.Role == "admin" {
		var admin entity.Admin
		if err := database.DB.Where("account_id = ?", user.Account_ID).First(&admin).Error; err == nil {
			// ถ้าเป็น admin เพิ่มข้อมูลเฉพาะของ admin
			userResp.AdminDetails = &response.AdminDetails{
				Admin_ID: admin.Admin_ID,
				Phone:    admin.Phone,
			}
		}
	}

	// ส่งข้อมูลกลับไป
	return ctx.Status(200).JSON(userResp)
}

func ManageAllUser(ctx fiber.Ctx) error {
	// ตรวจสอบสิทธิ์จาก JWT token
	claims := middleware.GetTokenClaims(ctx)
	adminRole := claims["role"].(string)
	adminUsername := claims["username"].(string)

	// ตรวจสอบว่ามีสิทธิ์เข้าถึงฟังก์ชันนี้หรือไม่
	if adminRole != "superadmin" && adminRole != "admin" {
		return ctx.Status(403).JSON(fiber.Map{
			"message": "ไม่มีสิทธิ์เข้าถึงฟังก์ชันนี้",
		})
	}

	// รับข้อมูลจาก request body
	request := struct {
		Account_ID uint   `json:"account_id" validate:"required"`
		Action     string `json:"action" validate:"required,oneof=change_status delete"`
		Status     string `json:"status" validate:"omitempty,oneof=Active Suspended"`
	}{}

	if err := ctx.Bind().Body(&request); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "ข้อมูลคำขอไม่ถูกต้อง: " + err.Error(),
		})
	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := validate.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// หาบัญชีผู้ใช้ที่ต้องการจัดการ
	var targetAccount entity.Account
	if err := database.DB.First(&targetAccount, request.Account_ID).Error; err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "ไม่พบบัญชีผู้ใช้ที่ต้องการจัดการ",
		})
	}

	// ตรวจสอบสิทธิ์ในการจัดการบัญชี
	if adminRole == "admin" {
		// Admin สามารถจัดการได้เฉพาะ user และ provider
		if targetAccount.Role == "admin" || targetAccount.Role == "superadmin" {
			return ctx.Status(403).JSON(fiber.Map{
				"message": "ไม่มีสิทธิ์จัดการบัญชีผู้ดูแลระบบหรือผู้ดูแลระบบสูงสุด",
			})
		}
	}

	// ดำเนินการตามคำสั่ง
	if request.Action == "change_status" {
		if request.Status == "" {
			return ctx.Status(400).JSON(fiber.Map{
				"message": "ต้องระบุสถานะที่ต้องการเปลี่ยน",
			})
		}

		// เปลี่ยนสถานะผู้ใช้
		if err := database.DB.Model(&targetAccount).Update("status", request.Status).Error; err != nil {
			return ctx.Status(409).JSON(fiber.Map{
				"message": "เกิดข้อผิดพลาดในการเปลี่ยนสถานะบัญชี: " + err.Error(),
			})
		}

		return ctx.Status(200).JSON(fiber.Map{
			"message":    "เปลี่ยนสถานะบัญชีผู้ใช้เรียบร้อยแล้ว",
			"account_id": targetAccount.Account_ID,
			"username":   targetAccount.Username,
			"status":     request.Status,
			"managed_by": adminUsername,
		})

	} else if request.Action == "delete" {
		// เริ่ม Transaction เพื่อลบข้อมูลที่เกี่ยวข้อง
		tx := database.DB.Begin()
		if tx.Error != nil {
			return ctx.Status(409).JSON(fiber.Map{
				"message": "ไม่สามารถเริ่ม Transaction ได้",
			})
		}

		// ลบข้อมูลตาม Role ของผู้ใช้
		if targetAccount.Role == "provider" {
			// ลบข้อมูล Provider
			if err := tx.Where("account_id = ?", targetAccount.Account_ID).Delete(&entity.Provider{}).Error; err != nil {
				tx.Rollback()
				return ctx.Status(409).JSON(fiber.Map{
					"message": "ไม่สามารถลบข้อมูล Provider ได้: " + err.Error(),
				})
			}
		} else if targetAccount.Role == "admin" {
			// ลบข้อมูล Admin
			if err := tx.Where("account_id = ?", targetAccount.Account_ID).Delete(&entity.Admin{}).Error; err != nil {
				tx.Rollback()
				return ctx.Status(409).JSON(fiber.Map{
					"message": "ไม่สามารถลบข้อมูล Admin ได้: " + err.Error(),
				})
			}
		}

		// ลบบัญชีผู้ใช้หลัก
		if err := tx.Delete(&targetAccount).Error; err != nil {
			tx.Rollback()
			return ctx.Status(409).JSON(fiber.Map{
				"message": "ไม่สามารถลบบัญชีผู้ใช้ได้: " + err.Error(),
			})
		}

		// Commit Transaction
		if err := tx.Commit().Error; err != nil {
			return ctx.Status(409).JSON(fiber.Map{
				"message": "เกิดข้อผิดพลาดในการลบบัญชีผู้ใช้: " + err.Error(),
			})
		}

		return ctx.Status(200).JSON(fiber.Map{
			"message":            "ลบบัญชีผู้ใช้เรียบร้อยแล้ว",
			"deleted_account_id": request.Account_ID,
			"deleted_username":   targetAccount.Username,
			"deleted_role":       targetAccount.Role,
			"managed_by":         adminUsername,
		})
	}

	return ctx.Status(400).JSON(fiber.Map{
		"message": "การดำเนินการไม่ถูกต้อง",
	})
}
