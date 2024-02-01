package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/config"
	"github.com/lengocson131002/go-clean/data/repo"
	"github.com/lengocson131002/go-clean/internal/interfaces"
	"github.com/lengocson131002/go-clean/internal/usecase"
	"github.com/lengocson131002/go-clean/pkg/database"
	"github.com/lengocson131002/go-clean/pkg/env"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"github.com/lengocson131002/go-clean/presenter/http/v1"
	"github.com/lengocson131002/go-clean/presenter/http/v1/controller"
	"github.com/lengocson131002/go-clean/presenter/http/v1/middleware"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(func() *env.ConfigFile {
			var file env.ConfigFile = ".env"
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
		fx.Provide(fx.Annotate(usecase.NewUserUseCase, fx.As(new(usecase.UserUseCase)))),
		fx.Provide(controller.NewUserController),
		fx.Provide(middleware.NewAuthMiddleware),
		fx.Provide(http.NewServer),
		fx.Invoke(run),
	).Run()
}

func run(lc fx.Lifecycle, app *fiber.App, log logger.Logger, conf *config.ServerConfig) {
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
