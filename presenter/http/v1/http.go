package http

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/lengocson131002/go-clean/config"
	"github.com/lengocson131002/go-clean/docs"
	"github.com/lengocson131002/go-clean/presenter/http/v1/controller"
	"github.com/lengocson131002/go-clean/presenter/http/v1/errors"
	"github.com/lengocson131002/go-clean/presenter/http/v1/middleware"
)

// @title  CLEAN ARCHITECTURE DEMO
// @version 1.0
// @description CLEAN ARCHITECTURE DEMO
// @termsOfService http://swagger.io/terms/
// @contact.name LNS
// @contact.email leson131002@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func NewServer(cfg *config.ServerConfig, userController *controller.UserController, authMiddleware *middleware.AuthMiddleware) *fiber.App {

	// middlewares
	app := fiber.New(fiber.Config{
		ErrorHandler: errors.CustomErrorHandler,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	// fiber log
	app.Use(fiberLog.New(fiberLog.Config{
		Next:         nil,
		Done:         nil,
		Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat:   "15:04:05",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
	}))

	// routes
	setSwagger(cfg)
	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	userRoute := v1.Group("/users")
	userRoute.Post("/register", userController.Register)
	userRoute.Get("/login", userController.Login)
	userRoute.Delete("/logout", authMiddleware.Handle, userController.Logout)
	userRoute.Get("/me", authMiddleware.Handle, userController.Current)
	userRoute.Put("", authMiddleware.Handle, userController.Update)

	return app
}

type Router struct {
	config         *config.ServerConfig
	Root           *fiber.App
	UserController *controller.UserController
	AuthMiddleware middleware.AuthMiddleware
}

func setSwagger(s *config.ServerConfig) {
	docs.SwaggerInfo.Title = "SWAGGER DOCUEMENTS"
	docs.SwaggerInfo.Description = "This is a go clean architecture example."
	docs.SwaggerInfo.Version = s.AppVersion
	docs.SwaggerInfo.Host = s.BaseURI
	docs.SwaggerInfo.BasePath = "/"
}
