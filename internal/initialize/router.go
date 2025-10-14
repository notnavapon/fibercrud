package initialize

import (
	"clean/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Router(deps *Dependencies, jwtSecret string) *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/users")
	api.Post("/register", deps.UserHandler.CreateUser)
	api.Post("/login", deps.UserHandler.LoginUser)
	api.Post("/logout", deps.UserHandler.LogoutUser)
	api.Get("", middleware.JWTProtected([]byte(jwtSecret)), deps.UserHandler.GetAllUser)
	api.Get("/:id", middleware.JWTProtected([]byte(jwtSecret)), deps.UserHandler.GetByID)
	api.Delete("/:id", middleware.JWTProtected([]byte(jwtSecret)), deps.UserHandler.Delete)
	api.Put("", middleware.JWTProtected([]byte(jwtSecret)), deps.UserHandler.Update)
	api.Patch("", middleware.JWTProtected([]byte(jwtSecret)), deps.UserHandler.ChangePassword)

	return app
}
