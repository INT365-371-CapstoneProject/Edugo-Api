package handler

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/middleware"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
	"gorm.io/gorm"
)

// กำหนดค่าเริ่มต้นสำหรับการตรวจสอบและแปลภาษา
func init() {
	enLocale := en.New()
	uni = ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")
	validate = validator.New()
	enTranslations.RegisterDefaultTranslations(validate, trans)

	// เพิ่มการแปล custom error messages
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	validate.RegisterTranslation("eqfield", trans, func(ut ut.Translator) error {
		return ut.Add("eqfield", "Confirm password must match password", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("eqfield")
		return t
	})
}

func CreateUser(ctx fiber.Ctx) error {
	user := new(request.UserRequest)
	if err := ctx.Bind().Body(user); err != nil {
		return utils.HandleError(ctx, 400, err.Error())
	}

	// Validate request
	if err := validate.Struct(user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// ส่งคืน error message แรกที่พบ
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// Check if email already exists
	var emailExists entity.Account
	result := database.DB.Where("email = ?", user.Email).First(&emailExists)
	if result.RowsAffected > 0 {
		return utils.HandleError(ctx, 400, "Email already exists")
	}

	// Check if username already exists
	var usernameExists entity.Account
	result = database.DB.Where("username = ?", user.Username).First(&usernameExists)
	if result.RowsAffected > 0 {
		return utils.HandleError(ctx, 400, "Username already exists")
	}

	// Begin transaction
	tx := database.DB.Begin()

	// create new user account
	newUser := entity.Account{
		Username:   user.Username,
		Email:      user.Email,
		FirstName:  &user.FirstName,
		LastName:   &user.LastName,
		Last_Login: nil,
		Role:       "user",
		Status:     "Active",
	}

	// Hashing password
	hashedPassword, err := utils.HashingPassword(user.Password)
	if err != nil {
		tx.Rollback()
		return utils.HandleError(ctx, 500, err.Error())
	}
	newUser.Password = hashedPassword

	// Insert account
	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return utils.HandleError(ctx, 500, err.Error())
	}

	// Create initial user record
	userRecord := entity.Users{
		Account_ID:      newUser.Account_ID,
		Education_level: nil, // กำหนดให้เป็น nil เพื่อให้เป็น NULL ในฐานข้อมูล
	}

	if err := tx.Create(&userRecord).Error; err != nil {
		tx.Rollback()
		return utils.HandleError(ctx, 500, err.Error())
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return utils.HandleError(ctx, 500, err.Error())
	}

	// สร้างข้อมูล Response
	userResponse := response.UserResponse{
		Account_ID: newUser.Account_ID,
		Username:   newUser.Username,
		Email:      newUser.Email,
		FirstName:  newUser.FirstName,
		LastName:   newUser.LastName,
		Create_On:  newUser.Create_On,
		Last_Login: newUser.Last_Login,
		Update_On:  newUser.Update_On,
		Role:       newUser.Role,
		Status:     newUser.Status,
	}

	return ctx.Status(201).JSON(userResponse)
}

func DeleteUser(ctx fiber.Ctx) error {
	claims := middleware.GetTokenClaims(ctx)
	username := claims["username"].(string)
	accountId := ctx.Params("id")

	var account entity.Account
	if err := database.DB.Where("username = ?", username).First(&account).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	var targetAccount entity.Account
	if err := database.DB.Where("account_id = ?", accountId).First(&targetAccount).Error; err != nil {
		return handleError(ctx, 404, "Target account not found")
	}

	if targetAccount.Account_ID != account.Account_ID {
		return handleError(ctx, 403, "Forbidden")
	}

	// แล้วค่อยลบตัว account
	if err := database.DB.Delete(&targetAccount).Error; err != nil {
		return handleError(ctx, 400, "Failed to delete account")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}

func CreateUserQuestion(ctx fiber.Ctx) error {
	// Get account_id from JWT token
	claims := middleware.GetTokenClaims(ctx)
	accountID := uint(claims["account_id"].(float64))

	// Check if user already answered questions
	var existingUser entity.Users
	if err := database.DB.Where("account_id = ? AND education_level IS NOT NULL", accountID).First(&existingUser).Error; err == nil {
		return utils.HandleError(ctx, 400, "User has already answered questions. Please use update endpoint instead.")
	} else if err != gorm.ErrRecordNotFound {
		return utils.HandleError(ctx, 500, err.Error())
	}

	question := new(request.QuestionRequest)
	if err := ctx.Bind().Body(question); err != nil {
		return utils.HandleError(ctx, 400, err.Error())
	}

	// Validate request
	if err := validate.Struct(question); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// Begin transaction
	tx := database.DB.Begin()

	// Get existing user record (should exist from CreateUser)
	var user entity.Users
	if err := tx.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		tx.Rollback()
		return utils.HandleError(ctx, 404, "User record not found")
	}

	// Update user with education level
	user.Education_level = &question.Education_Level
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return utils.HandleError(ctx, 500, err.Error())
	}

	// Create answer_countries records
	for _, countryID := range question.Countries {
		answer := entity.AnswerCountries{
			User_ID:    user.User_ID,
			Country_ID: countryID,
		}
		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			return utils.HandleError(ctx, 500, err.Error())
		}
	}

	// Create answer_categories records
	for _, categoryID := range question.Categories {
		answer := entity.AnswerCategories{
			User_ID:     user.User_ID,
			Category_ID: categoryID,
		}
		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			return utils.HandleError(ctx, 500, err.Error())
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return utils.HandleError(ctx, 500, err.Error())
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "User questions saved successfully",
	})
}

func UpdateUserQuestion(ctx fiber.Ctx) error {
	// Get account_id from JWT token
	claims := middleware.GetTokenClaims(ctx)
	accountID := uint(claims["account_id"].(float64))

	// Parse request body
	question := new(request.QuestionRequest)
	if err := ctx.Bind().Body(question); err != nil {
		return utils.HandleError(ctx, 400, err.Error())
	}

	// Validate request
	if err := validate.Struct(question); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.HandleError(ctx, 400, validationErrors[0].Translate(trans))
	}

	// Begin transaction
	tx := database.DB.Begin()

	// Get user record
	var user entity.Users
	if err := tx.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(ctx, 404, "User not found")
		}
		return utils.HandleError(ctx, 409, err.Error())
	}

	// Update education level
	user.Education_level = &question.Education_Level
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return utils.HandleError(ctx, 409, err.Error())
	}

	// Delete existing answer_countries
	if err := tx.Where("user_id = ?", user.User_ID).Delete(&entity.AnswerCountries{}).Error; err != nil {
		tx.Rollback()
		return utils.HandleError(ctx, 409, err.Error())
	}
	if err := tx.Where("user_id = ?", user.User_ID).Delete(&entity.AnswerCategories{}).Error; err != nil {
		tx.Rollback()
		return utils.HandleError(ctx, 409, err.Error())
	}

	// Create new answer_countries records
	for _, countryID := range question.Countries {
		answer := entity.AnswerCountries{
			User_ID:    user.User_ID,
			Country_ID: countryID,
		}
		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			return utils.HandleError(ctx, 409, err.Error())
		}
	}

	// Create new answer_categories records
	for _, categoryID := range question.Categories {
		answer := entity.AnswerCategories{
			User_ID:     user.User_ID,
			Category_ID: categoryID,
		}
		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			return utils.HandleError(ctx, 409, err.Error())
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return utils.HandleError(ctx, 409, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"message": "User questions updated successfully",
	})
}

func GetUserQuestions(ctx fiber.Ctx) error {
	// Get account_id from JWT token
	claims := middleware.GetTokenClaims(ctx)
	accountID := uint(claims["account_id"].(float64))

	var user entity.Users
	if err := database.DB.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(ctx, 404, "User questions not found")
		}
		return utils.HandleError(ctx, 409, err.Error())
	}

	var answers []entity.AnswerCountries
	if err := database.DB.Preload("Country").Where("user_id = ?", user.User_ID).Find(&answers).Error; err != nil {
		return utils.HandleError(ctx, 409, err.Error())
	}

	var categoryAnswers []entity.AnswerCategories
	if err := database.DB.Preload("Category").Where("user_id = ?", user.User_ID).Find(&categoryAnswers).Error; err != nil {
		return utils.HandleError(ctx, 409, err.Error())
	}

	// Create response
	countries := make([]map[string]interface{}, len(answers))
	for i, answer := range answers {
		countries[i] = map[string]interface{}{
			"country_id": answer.Country_ID,
			"name":       answer.Country.Name,
		}
	}

	categories := make([]map[string]interface{}, len(categoryAnswers))
	for i, answer := range categoryAnswers {
		categories[i] = map[string]interface{}{
			"category_id": answer.Category_ID,
			"name":        answer.Category.Name,
		}
	}

	return ctx.JSON(fiber.Map{
		"education_level": user.Education_level,
		"countries":       countries,
		"categories":      categories,
	})
}
