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
		})
	}

	// ChECK AVAILABLE User
	var user entity.User
	err := database.DB.Preload("Role").First(&user, "email = ? ", loginRequest.Email).Error
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Invalid credentials",
		})
	}

	// CHECK VALIDATION Password
	isValid := utils.CheckPasswordHash(loginRequest.Password, user.Password)
	if !isValid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Invalid credentials",
		})
	}

	// GENERATE JWT
	claims := jwt.MapClaims{}
	claims["full_name"] = user.FullName
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
	})
}
