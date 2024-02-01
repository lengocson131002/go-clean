package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/api"
	"github.com/lengocson131002/go-clean/api/controller"
	"github.com/lengocson131002/go-clean/api/middleware"
	"github.com/lengocson131002/go-clean/bootstrap"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/repository"
	"github.com/lengocson131002/go-clean/usecase"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(bootstrap.GetConfigure),
		fx.Provide(bootstrap.GetLogger),
		fx.Provide(bootstrap.GetValidator),
		fx.Provide(bootstrap.GetDatabaseConfig),
		fx.Provide(bootstrap.GetServerConfig),
		fx.Provide(bootstrap.GetDatabaseConnector),
		fx.Provide(bootstrap.GetDatabase),
		fx.Provide(fx.Annotate(repository.NewUserRepository, fx.As(new(domain.UserRepository)))),
		fx.Provide(fx.Annotate(usecase.NewUserUseCase, fx.As(new(domain.UserUseCase)))),
		fx.Provide(controller.NewUserController),
		fx.Provide(middleware.NewAuthMiddleware),
		fx.Provide(api.NewHttpServer),
		fx.Invoke(run),
	).Run()
}

func run(lc fx.Lifecycle, app *fiber.App, log logger.Logger, conf *bootstrap.ServerConfig) {
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
