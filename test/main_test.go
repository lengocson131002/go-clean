package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/config"
	"github.com/lengocson131002/go-clean/data/repo"
	"github.com/lengocson131002/go-clean/internal/interfaces"
	"github.com/lengocson131002/go-clean/internal/model"
	"github.com/lengocson131002/go-clean/internal/usecase"
	"github.com/lengocson131002/go-clean/pkg/database"
	"github.com/lengocson131002/go-clean/pkg/env"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"github.com/lengocson131002/go-clean/presenter/http/v1"
	"github.com/lengocson131002/go-clean/presenter/http/v1/controller"
	"github.com/lengocson131002/go-clean/presenter/http/v1/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestDatabaseTransaction(t *testing.T) {
	fx.New(
		fx.Provide(func() *env.ConfigFile {
			var file env.ConfigFile = "test.env"
			return &file
		}),
		fx.Provide(fx.Annotate(env.NewViperConfig, fx.As(new(env.Configure)))),
		fx.Provide(fx.Annotate(logger.NewLogrus, fx.As(new(logger.Logger)))),
		fx.Provide(fx.Annotate(validation.NewGpValidator, fx.As(new(validation.Validator)))),
		fx.Provide(config.GetDatabaseConfig),
		fx.Provide(config.GetServerConfig),
		fx.Provide(config.GetDatabaseSqlx),
		fx.Provide(database.GetSqlxGdbc),
		fx.Provide(fx.Annotate(repo.NewUserRepository, fx.As(new(interfaces.UserRepositoryInterface)))),
		fx.Provide(usecase.NewUserUseCase),
		fx.Provide(controller.NewUserController),
		fx.Provide(middleware.NewAuthMiddleware),
		fx.Provide(http.NewServer),
		// fx.Invoke(start),
		fx.Provide(func() *testing.T {
			return t
		}),
		fx.Invoke(runTest),
	).Run()
}

func start(lc fx.Lifecycle, app *fiber.App, log logger.Logger, conf *config.ServerConfig) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(fmt.Sprintf(":%v", conf.Port)); err != nil {
					log.Error("server terminated unexpectedly")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("gracefully shutting down server")
			if err := app.Shutdown(); err != nil {
				log.Error("error occurred while gracefully shutting down server")
				return err
			}
			log.Info("graceful server shut down completed")
			return nil
		},
	})
}

func runTest(shutdowner fx.Shutdowner, t *testing.T, app *fiber.App, log logger.Logger) {
	type Request struct {
		id       string
		password string
		name     string
	}

	requests := make([]model.RegisterUserRequest, 0)

	for i := 1; i <= 1000; i++ {
		r := model.RegisterUserRequest{
			ID:       fmt.Sprintf("%v", i),
			Password: "12345678",
			Name:     "test",
		}
		requests = append(requests, r)
	}

	for _, reqData := range requests {
		byteData, _ := json.Marshal(&reqData)
		log.Info("Request: %s", byteData)
		req := httptest.NewRequest("POST", "/api/v1/users/register", bytes.NewReader(byteData))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 100)

		assert.Equal(t, 200, resp.StatusCode)
	}

	app.Shutdown()
	shutdowner.Shutdown()
}
