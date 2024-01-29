package http

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/lengocson131002/go-clean/config"
	"github.com/lengocson131002/go-clean/data/repo"
	"github.com/lengocson131002/go-clean/docs"
	"github.com/lengocson131002/go-clean/internal/usecase"
	"github.com/lengocson131002/go-clean/pkg/database"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"github.com/lengocson131002/go-clean/presenter/http/v1/controller"
	"github.com/lengocson131002/go-clean/presenter/http/v1/errors"
	"github.com/lengocson131002/go-clean/presenter/http/v1/middleware"
	"github.com/lengocson131002/go-clean/presenter/http/v1/route"
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
func RunServer(cfg *config.BootstrapConfig) error {
	dbConfig := config.GetDatabaseConfig(cfg.Config)
	db, err := config.GetDatabaseSqlx(dbConfig)
	if err != nil {
		cfg.Logger.Fatal("failed to connect to database %s", err.Error())
	}

	serverConfig := config.GetServerConfig(cfg.Config)

	// repositories
	userRepo := repo.NewUserRepository(database.GetSqlxGdbc(db))

	// validator
	validator := validation.NewGpValidator()

	// usercases
	userUseCase := usecase.NewUserUseCase(cfg.Logger, validator, userRepo)

	// controllers
	userController := controller.NewUserController(userUseCase, cfg.Logger)

	// middlewares
	authMiddleware := middleware.NewAuth(userUseCase, cfg.Logger)

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

	// swagger

	setSwagger(serverConfig)
	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	userRoute := route.RouteConfig{
		Root:           &v1,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}

	userRoute.Setup()

	// start the server
	err = app.Listen(fmt.Sprintf(":%d", serverConfig.Port))
	if err != nil {
		cfg.Logger.Fatal("Failed to start server: %v", err)
		return err
	}

	return nil
}

func setSwagger(s *config.ServerConfig) {
	docs.SwaggerInfo.Title = "SWAGGER DOCUEMENTS"
	docs.SwaggerInfo.Description = "This is a go clean architecture example."
	docs.SwaggerInfo.Version = s.AppVersion
	docs.SwaggerInfo.Host = s.BaseURI
	docs.SwaggerInfo.BasePath = "/"
}
