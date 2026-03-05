package handler

import (
	"go-cariir-api/database"
	"go-cariir-api/model/entity"
	"go-cariir-api/model/request"
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
	}

	return ctx.JSON(fiber.Map{"data": users})
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
		return ctx.Status(400).JSON(fiber.Map{"message": "failed", "errors": errValidate.Error()})
	}

	newUser := entity.User{
		FullName: user.FullName,
		Email:    user.Email,
		RoleID:   user.RoleID,
	}

	hashedPassword, err := utils.HashingPassword(user.Password)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal server error", "errors": err.Error()})
	}

	newUser.Password = hashedPassword

	errCreate := database.DB.Create(&newUser).Error

	if errCreate != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": errCreate.Error()})
	}

	return ctx.JSON(fiber.Map{"data": newUser, "message": "user created"})
}

func UserHandlerGetById(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")

	var user entity.User

	result := database.DB.First(&user, "id = ? ", userId).Error
	if result != nil {
		return ctx.Status(404).JSON(fiber.Map{"message": "user not found"})
	}

	return ctx.JSON(fiber.Map{"data": user, "message": "success"})
}

func UserHandlerUpdate(ctx *fiber.Ctx) error {
	userRequest := new(request.UserUpdateRequest)

	if err := ctx.BodyParser(userRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "bad request"})
	}

	validate := validator.New()
	errValidate := validate.Struct(userRequest)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "failed", "errors": errValidate.Error()})
	}

	userId := ctx.Params("id")

	var user entity.User

	result := database.DB.First(&user, "id = ? ", userId).Error
	if result != nil {
		return ctx.Status(404).JSON(fiber.Map{"message": "user not found"})
	}

	user.FullName = userRequest.FullName
	user.IsActive = userRequest.IsActive

	errUpdate := database.DB.Save(&user).Error
	if errUpdate != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": errUpdate.Error()})
	}

	return ctx.JSON(fiber.Map{"data": user, "message": "success"})
}

func UserHandlerDelete(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")

	var user entity.User

	result := database.DB.Debug().First(&user, "id = ? ", userId).Error
	if result != nil {
		return ctx.Status(404).JSON(fiber.Map{"message": "user not found"})
	}

	errDelete := database.DB.Debug().Delete(&user).Error

	if errDelete != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": "internal server error"})
	}

	return ctx.JSON(fiber.Map{"message": "user deleted"})
}
