package middleware

import (
	"go-cariir-api/database"
	"go-cariir-api/model/response"
	"go-cariir-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(ctx *fiber.Ctx) error {
	authorization := ctx.Get("Authorization")

	if authorization == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Unauthenticated",
			Code:    fiber.StatusUnauthorized,
		})
	}

	if len(authorization) < 7 || authorization[:7] != "Bearer " {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Invalid authorization header format",
			Code:    fiber.StatusUnauthorized,
		})
	}

	token := authorization[7:]

	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Unauthenticated",
			Code:    fiber.StatusUnauthorized,
		})
	}

	ctx.Locals("user", claims)
	return ctx.Next()
}

func PermissionGuard(permission string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userClaims := ctx.Locals("user").(jwt.MapClaims)
		roleAlias := userClaims["role"].(string)

		if roleAlias == "admin" {
			return ctx.Next()
		}

		var count int64
		err := database.DB.Table("role_permissions").
			Joins("join roles on roles.id = role_permissions.role_id").
			Joins("join permissions on permissions.id = role_permissions.permission_id").
			Where("roles.alias = ? and permissions.alias = ?", roleAlias, permission).
			Count(&count).Error

		if err != nil || count == 0 {
			return ctx.Status(fiber.StatusForbidden).JSON(response.GenericResponse{
				Status:  "error",
				Message: "Forbidden: You don't have permission to access " + permission,
				Code:    fiber.StatusForbidden,
			})
		}

		return ctx.Next()
	}
}

func PermissionCreate(ctx *fiber.Ctx) error {
	return ctx.Next()
}
