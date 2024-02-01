package api

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/lengocson131002/go-clean/api/controller"
	"github.com/lengocson131002/go-clean/api/handler"
	"github.com/lengocson131002/go-clean/api/middleware"
	"github.com/lengocson131002/go-clean/api/route"
	"github.com/lengocson131002/go-clean/bootstrap"
	"github.com/lengocson131002/go-clean/docs"
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
func NewHttpServer(cfg *bootstrap.ServerConfig, userController *controller.UserController, authMiddleware *middleware.AuthMiddleware) *fiber.App {

	// middlewares
	app := fiber.New(fiber.Config{
		ErrorHandler: handler.CustomErrorHandler,
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

	// Register routes
	route.RegisterUserRoute(&v1, userController, authMiddleware)

	return app
}

type Router struct {
	config         *bootstrap.ServerConfig
	Root           *fiber.App
	UserController *controller.UserController
	AuthMiddleware middleware.AuthMiddleware
}

func setSwagger(s *bootstrap.ServerConfig) {
	docs.SwaggerInfo.Title = "SWAGGER DOCUEMENTS"
	docs.SwaggerInfo.Description = "This is a go clean architecture example."
	docs.SwaggerInfo.Version = s.AppVersion
	docs.SwaggerInfo.Host = s.BaseURI
	docs.SwaggerInfo.BasePath = "/"
}
