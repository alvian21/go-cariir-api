package route

import (
	"go-cariir-api/config"
	"go-cariir-api/handler"
	"go-cariir-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {

	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Application is running",
			"status":  "ok",
		})
	})

	r.Static("/public", config.ProjectRootPath+"/public/asset")

	// Auth
	auth := r.Group("/auth")
	auth.Post("/login", handler.LoginHandler)
	auth.Post("/register", handler.RegisterHandler)

	// User
	user := r.Group("/user", middleware.Auth)
	user.Get("/", middleware.PermissionGuard("user.read"), handler.UserHandlerGetAll)
	user.Post("/", middleware.PermissionGuard("user.create"), handler.UserHandlerCreate)
	user.Get("/:id", middleware.PermissionGuard("user.read"), handler.UserHandlerGetById)
	user.Put("/:id", middleware.PermissionGuard("user.update"), handler.UserHandlerUpdate)
	user.Delete("/:id", middleware.PermissionGuard("user.delete"), handler.UserHandlerDelete)
}
