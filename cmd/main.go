package main

import (
	"context"
	"time"

	"github.com/lengocson131002/go-clean/bootstrap"
	"github.com/lengocson131002/go-clean/infras/data"
	"github.com/lengocson131002/go-clean/infras/outbound"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/xslt"
	gprc "github.com/lengocson131002/go-clean/presentation/grpc"
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
	fx.Provide(bootstrap.NewErrorHandlingBehavior),

	fx.Provide(http.NewHttpServer),
	fx.Provide(gprc.NewGrpcServer),
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

func run(lc fx.Lifecycle, http *http.HttpServer, grpc *gprc.GrpcServer, log logger.Logger, conf *bootstrap.ServerConfig, shutdowner fx.Shutdowner) {
	gCtx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			errChan := make(chan error)

			// start HTTP server
			go func() {
				if err := http.Start(gCtx); err != nil {
					log.Fatal("Failed to start HTTP server: %s", err)
					errChan <- err
					cancel()
					shutdowner.Shutdown()
				}
			}()

			// start GRPC server
			go func() {
				if err := grpc.Start(gCtx); err != nil {
					log.Fatal("Failed to start GRPC server: %s", err)
					errChan <- err
					cancel()
					shutdowner.Shutdown()
				}
			}()

			select {
			case err := <-errChan:
				return err
			case <-time.After(100 * time.Millisecond):
				return nil
			}

		},
		OnStop: func(ctx context.Context) error {
			cancel()
			select {
			case <-time.After(100 * time.Millisecond):
				return nil
			}
		},
	})
}
