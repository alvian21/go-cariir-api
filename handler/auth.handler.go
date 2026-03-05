package handler

import (
	"go-cariir-api/database"
	"go-cariir-api/model/entity"
	"go-cariir-api/model/request"
	"go-cariir-api/model/response"
	"go-cariir-api/utils"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(ctx *fiber.Ctx) error {
	loginRequest := new(request.LoginRequest)

	if err := ctx.BodyParser(loginRequest); err != nil {
		log.Println(err)
		return err
	}

	validate := validator.New()
	errValidate := validate.Struct(loginRequest)
	if errValidate != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Validation failed",
			Errors:  errValidate.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	// ChECK AVAILABLE User
	var user entity.User
	err := database.DB.Preload("Role").First(&user, "email = ? ", loginRequest.Email).Error
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Invalid credentials",
			Code:    fiber.StatusUnauthorized,
		})
	}

	// CHECK VALIDATION Password
	isValid := utils.CheckPasswordHash(loginRequest.Password, user.Password)
	if !isValid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Invalid credentials",
			Code:    fiber.StatusUnauthorized,
		})
	}

	// GENERATE JWT
	claims := jwt.MapClaims{}
	claims["userId"] = user.ID
	claims["fullName"] = user.FullName
	claims["email"] = user.Email
	if user.Role != nil {
		claims["role"] = user.Role.Alias
	}
	claims["exp"] = time.Now().Add(time.Minute * 2).Unix()

	token, errGenerateToken := utils.GenerateToken(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Failed to generate token",
			Code:    fiber.StatusUnauthorized,
		})
	}

	userResponse := response.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		IsActive:  user.IsActive,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	loginResponse := response.LoginResponse{
		AccessToken: token,
		ExpiresAt:   claims["exp"].(int64),
		User:        userResponse,
	}

	return ctx.JSON(response.GenericResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    loginResponse,
		Code:    fiber.StatusOK,
	})
}

func RegisterHandler(ctx *fiber.Ctx) error {
	registerRequest := new(request.RegisterRequest)

	if err := ctx.BodyParser(registerRequest); err != nil {
		log.Println(err)
		return err
	}

	validate := validator.New()
	errValidate := validate.Struct(registerRequest)
	if errValidate != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Validation failed",
			Errors:  errValidate.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	// CHECK AVAILABLE User
	var user entity.User
	err := database.DB.First(&user, "email = ? ", registerRequest.Email).Error
	if err == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Email already exists",
			Code:    fiber.StatusBadRequest,
		})
	}

	// CREATE USER
	hashedPassword, errHash := utils.HashingPassword(registerRequest.Password)
	if errHash != nil {
		log.Println(errHash)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Failed to hash password",
			Code:    fiber.StatusInternalServerError,
		})
	}

	var role entity.Role
	err = database.DB.First(&role, "alias = ?", "user").Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Failed to find role",
			Code:    fiber.StatusInternalServerError,
		})
	}

	user = entity.User{
		FullName: registerRequest.FullName,
		Email:    registerRequest.Email,
		Password: hashedPassword,
		IsActive: true,
		RoleID:   &role.ID,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Failed to create user",
			Code:    fiber.StatusInternalServerError,
		})
	}

	userResponse := response.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		IsActive:  user.IsActive,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return ctx.JSON(response.GenericResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data:    userResponse,
		Code:    fiber.StatusOK,
	})
}
