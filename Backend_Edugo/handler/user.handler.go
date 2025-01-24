package handler

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v3"
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"github.com/tk-neng/demo-go-fiber/request"
	"github.com/tk-neng/demo-go-fiber/response"
	"github.com/tk-neng/demo-go-fiber/utils"
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

func GetAllUser(ctx fiber.Ctx) error {
	var users []entity.Account
	result := database.DB.Where("role = ?", "user").Find(&users)
	if result.Error != nil {
		return utils.HandleError(ctx, 404, result.Error.Error())
	}
	
	var userResponse []response.UserResponse
	for _, user := range users {
		userResponse = append(userResponse, response.UserResponse{
			Account_ID: user.Account_ID,
			Username: user.Username,
			Email: user.Email,
			Create_On: user.Create_On,
			Last_Login: user.Last_Login,
			Update_On: user.Update_On,
			Role: user.Role,
		})
	}
	return ctx.Status(200).JSON(userResponse)
}

func GetUserByID(ctx fiber.Ctx) error {
	// Get parameter value
	accountId := ctx.Params("id")
	var user entity.Account
	result := database.DB.Where("account_id = ? AND role = ?", accountId, "user").First(&user)
	if result.Error != nil {
		return utils.HandleError(ctx, 404, "User not found")
	}

	userResponse := response.UserResponse{
		Account_ID: user.Account_ID,
		Username: user.Username,
		Email: user.Email,
		Create_On: user.Create_On,
		Last_Login: user.Last_Login,
		Update_On: user.Update_On,
		Role: user.Role,
	}

	return ctx.Status(200).JSON(userResponse)
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

	// create new user
	newUser := entity.Account{
		Username: user.Username,
		Email: user.Email,
		Last_Login: nil,
		Role: "user",
	}

	// Hashing password
	hashedPassword, err := utils.HashingPassword(user.Password)
	if err != nil {
		return utils.HandleError(ctx, 500, err.Error())
	}
	newUser.Password = hashedPassword

	// Insert to database
	errCreateUser := database.DB.Create(&newUser).Error
	if errCreateUser != nil {
		return utils.HandleError(ctx, 500, errCreateUser.Error())
	}


	// สร้างข้อมูล Response
	userResponse := response.UserResponse{
		Account_ID: newUser.Account_ID,
		Username: newUser.Username,
		Email: newUser.Email,
		Create_On: newUser.Create_On,
		Last_Login: newUser.Last_Login,
		Update_On: newUser.Update_On,
		Role: newUser.Role,
	}

	return ctx.Status(201).JSON(userResponse)
}

