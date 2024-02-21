package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/bootstrap"
	"github.com/lengocson131002/go-clean/infras/data"
	"github.com/lengocson131002/go-clean/infras/outbound"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/xslt"
	"github.com/lengocson131002/go-clean/presentation/http"
	"github.com/lengocson131002/go-clean/presentation/http/controller"
	"github.com/lengocson131002/go-clean/presentation/http/middleware"
	"github.com/lengocson131002/go-clean/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module("main",
	fx.Provide(bootstrap.GetLogger),
	fx.Provide(bootstrap.GetConfigure),
	fx.Provide(bootstrap.GetServerConfig),
	fx.Provide(bootstrap.GetDatabaseConfig),
	fx.Provide(bootstrap.GetTracingConfig),
	fx.Provide(bootstrap.GetValidator),
	fx.Provide(bootstrap.GetDatabaseConnector),
	fx.Provide(bootstrap.GetUserDatabase),
	fx.Provide(bootstrap.GetTracer),
	fx.Provide(data.NewUserRepository),
	fx.Provide(controller.NewUserController),
	fx.Provide(middleware.NewAuthMiddleware),
	fx.Provide(bootstrap.NewHealthChecker),

	fx.Provide(usecase.NewVerifyUserHandler),
	fx.Provide(usecase.NewCreateUserHandler),
	fx.Provide(usecase.NewGetUserHandler),
	fx.Provide(usecase.NewLogoutUserHandler),
	fx.Provide(usecase.NewUpdateUserHandler),
	fx.Provide(usecase.NewLoginUserHandler),
	fx.Provide(usecase.NewOpenAccountHandler),

	fx.Provide(bootstrap.NewRequestLoggingBehavior),
	fx.Provide(bootstrap.NewTracingBehavior),
	fx.Provide(bootstrap.NewMetricBehavior),
	fx.Provide(http.NewHttpServer),
	fx.Provide(bootstrap.GetYugabyteConfig),
	fx.Provide(bootstrap.GetMasterDataDatabase),
	fx.Provide(data.NewMasterDataRepository),
	fx.Provide(bootstrap.NewMetricer),
	fx.Provide(bootstrap.GetT24MqConfig),
	fx.Provide(controller.NewT24AccountController),
	fx.Provide(xslt.NewDefaultXslt),
	fx.Provide(outbound.NewT24MqClient),
	fx.Invoke(bootstrap.RegisterPipeline),
)

func main() {
	// Dependencies injection using FX package
	fx.New(
		Module,
		fx.Invoke(run),
	).Run()
}

func run(lc fx.Lifecycle, app *fiber.App, log logger.Logger, conf *bootstrap.ServerConfig, shutdowner fx.Shutdowner) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(fmt.Sprintf(":%v", conf.Port)); err != nil {
					log.Error("server terminated unexpectedly")
					shutdowner.Shutdown()
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
