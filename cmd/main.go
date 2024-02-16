package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/bootstrap"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/infras/repository"
	"github.com/lengocson131002/go-clean/pkg/logger"
	api "github.com/lengocson131002/go-clean/presentation/http"
	"github.com/lengocson131002/go-clean/presentation/http/controller"
	"github.com/lengocson131002/go-clean/presentation/http/middleware"
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
		fx.Provide(bootstrap.GetTracingConfig),
		fx.Provide(bootstrap.GetTracer),
		fx.Provide(middleware.NewTracingMiddleware),
		fx.Provide(fx.Annotate(repository.NewUserRepository, fx.As(new(domain.UserRepository)))),
		fx.Provide(fx.Annotate(usecase.NewUserUseCase, fx.As(new(domain.UserUseCase)))),
		fx.Provide(controller.NewUserController),
		fx.Provide(middleware.NewAuthMiddleware),
		fx.Provide(fx.Annotate(bootstrap.NewHealthEndpoint, fx.As(new(bootstrap.HealthCheckerEndpoint)))),
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
