package handler

import (
	"go-cariir-api/database"
	"go-cariir-api/model/entity"
	"go-cariir-api/model/request"
	"go-cariir-api/model/response"
	"go-cariir-api/utils"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UserHandlerGetAll(ctx *fiber.Ctx) error {

	userInfo := ctx.Locals("user")
	log.Println("user info :", userInfo)

	var users []entity.User
	result := database.DB.Find(&users)

	if result.Error != nil {
		log.Println(result.Error)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Failed to fetch users",
			Code:    fiber.StatusInternalServerError,
		})
	}

	var usersResponse []response.UserResponse
	for _, u := range users {
		usersResponse = append(usersResponse, response.UserResponse{
			ID:        u.ID,
			FullName:  u.FullName,
			Email:     u.Email,
			IsActive:  u.IsActive,
			RoleID:    u.RoleID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return ctx.JSON(response.GenericResponse{
		Status:  "success",
		Message: "Success fetch users",
		Data:    usersResponse,
		Code:    fiber.StatusOK,
	})
}

func UserHandlerCreate(ctx *fiber.Ctx) error {
	user := new(request.UserCreateRequest)

	if err := ctx.BodyParser(user); err != nil {
		log.Println(err)
		return err
	}

	validate := validator.New()
	errValidate := validate.Struct(user)
	if errValidate != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Validation failed",
			Errors:  errValidate.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	newUser := entity.User{
		FullName: user.FullName,
		Email:    user.Email,
		RoleID:   user.RoleID,
	}

	hashedPassword, err := utils.HashingPassword(user.Password)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Failed to hash password",
			Code:    fiber.StatusInternalServerError,
		})
	}

	newUser.Password = hashedPassword

	errCreate := database.DB.Create(&newUser).Error

	if errCreate != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Failed to create user",
			Code:    fiber.StatusInternalServerError,
		})
	}

	userResponse := response.UserResponse{
		ID:        newUser.ID,
		FullName:  newUser.FullName,
		Email:     newUser.Email,
		IsActive:  newUser.IsActive,
		RoleID:    newUser.RoleID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return ctx.JSON(response.GenericResponse{
		Status:  "success",
		Message: "User created successfully",
		Data:    userResponse,
		Code:    fiber.StatusOK,
	})
}

func UserHandlerGetById(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")

	var user entity.User

	result := database.DB.First(&user, "id = ? ", userId).Error
	if result != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.GenericResponse{
			Status:  "error",
			Message: "User not found",
			Code:    fiber.StatusNotFound,
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
		Message: "Success fetch user",
		Data:    userResponse,
		Code:    fiber.StatusOK,
	})
}

func UserHandlerUpdate(ctx *fiber.Ctx) error {
	userRequest := new(request.UserUpdateRequest)

	if err := ctx.BodyParser(userRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "bad request"})
	}

	validate := validator.New()
	errValidate := validate.Struct(userRequest)
	if errValidate != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Validation failed",
			Errors:  errValidate.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	userId := ctx.Params("id")

	var user entity.User

	result := database.DB.First(&user, "id = ? ", userId).Error
	if result != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.GenericResponse{
			Status:  "error",
			Message: "User not found",
			Code:    fiber.StatusNotFound,
		})
	}

	user.FullName = userRequest.FullName
	user.IsActive = userRequest.IsActive

	errUpdate := database.DB.Save(&user).Error
	if errUpdate != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Failed to update user",
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
		Message: "User updated successfully",
		Data:    userResponse,
		Code:    fiber.StatusOK,
	})
}

func UserHandlerDelete(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")

	var user entity.User

	result := database.DB.Debug().First(&user, "id = ? ", userId).Error
	if result != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.GenericResponse{
			Status:  "error",
			Message: "User not found",
			Code:    fiber.StatusNotFound,
		})
	}

	errDelete := database.DB.Debug().Delete(&user).Error

	if errDelete != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Failed to delete user",
			Code:    fiber.StatusInternalServerError,
		})
	}

	return ctx.JSON(response.GenericResponse{
		Status:  "success",
		Message: "User deleted successfully",
		Code:    fiber.StatusOK,
	})
}
