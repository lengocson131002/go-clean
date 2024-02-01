package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/api/controller"
	"github.com/lengocson131002/go-clean/api/middleware"
)

func RegisterUserRoute(root *fiber.Router, userController *controller.UserController, authMiddleware *middleware.AuthMiddleware) {
	userRoute := (*root).Group("/users")

	userRoute.Post("/register", userController.Register)
	userRoute.Get("/login", userController.Login)
	userRoute.Delete("/logout", authMiddleware.Handle, userController.Logout)
	userRoute.Get("/me", authMiddleware.Handle, userController.Current)
	userRoute.Put("", authMiddleware.Handle, userController.Update)
}
