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

	// User
	user := r.Group("/user", middleware.Auth)
	user.Get("/", handler.UserHandlerGetAll)
	user.Get("/:id", handler.UserHandlerGetById)
	user.Post("/", handler.UserHandlerCreate)
	user.Put("/:id", handler.UserHandlerUpdate)
	user.Delete("/:id", handler.UserHandlerDelete)
}
