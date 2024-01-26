package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/presenter/http/v1/controller"
)

type RouteConfig struct {
	Root           *fiber.Router
	UserController *controller.UserController
	AuthMiddleware fiber.Handler
}

func (r *RouteConfig) Setup() {
	userRoute := (*r.Root).Group("/users")

	userRoute.Post("/register", r.UserController.Register)
	userRoute.Get("/login", r.UserController.Login)
	userRoute.Delete("/logout", r.AuthMiddleware, r.UserController.Logout)
	userRoute.Get("/me", r.AuthMiddleware, r.UserController.Current)

}
