package http

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/lengocson131002/go-clean/bootstrap"
	healthchecks "github.com/lengocson131002/go-clean/pkg/health"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/presentation/http/controller"
	"github.com/lengocson131002/go-clean/presentation/http/handler"
	"github.com/lengocson131002/go-clean/presentation/http/route"
)

type HttpServer struct {
	cfg               *bootstrap.ServerConfig
	logger            logger.Logger
	healhChecker      healthchecks.HealthChecker
	t24AccConntroller *controller.T24AccountController
}

// @title  CLEAN ARCHITECTURE DEMO
// @version 1.0
// @description CLEAN ARCHITECTURE DEMO
// @termsOfService http://swagger.io/terms/
// @contact.name LNS
// @contact.email leson131002@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func NewHttpServer(
	cfg *bootstrap.ServerConfig,
	logger logger.Logger,
	healhChecker healthchecks.HealthChecker,
	t24AccConntroller *controller.T24AccountController) *HttpServer {
	return &HttpServer{
		cfg:               cfg,
		logger:            logger,
		healhChecker:      healhChecker,
		t24AccConntroller: t24AccConntroller,
	}
}

func (s *HttpServer) Start(ctx context.Context) error {
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
		TimeFormat:   "2006-01-02 15:04:05",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	// health check endpoint
	app.Get("/liveliness", func(c *fiber.Ctx) error {
		result := s.healhChecker.LivenessCheck()
		if result.Status {
			return c.Status(fiber.StatusOK).JSON(result)
		}
		return c.Status(fiber.StatusServiceUnavailable).JSON(result)
	})

	app.Get("/readiness", func(c *fiber.Ctx) error {
		result := s.healhChecker.RedinessCheck()
		if result.Status {
			return c.Status(fiber.StatusOK).JSON(result)
		}
		return c.Status(fiber.StatusServiceUnavailable).JSON(result)
	})

	// tracing
	app.Use(otelfiber.Middleware())

	// metrics endpoint
	prometheus := fiberprometheus.New(s.cfg.Name)
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Register routes
	route.RegisterT24Route(&v1, s.t24AccConntroller)

	go func() {
		defer func(ctx context.Context) {
			if err := app.Shutdown(); err != nil {
				s.logger.Errorf(ctx, "Failed to shutdown http server: %v", err)
			}
			s.logger.Info(ctx, "Stop HTTP Server")
		}(ctx)

		<-ctx.Done()
	}()

	hPort := s.cfg.HttpPort
	s.logger.Infof(ctx, "Start HTTP server at port: %v", hPort)
	if err := app.Listen(fmt.Sprintf(":%v", hPort)); err != nil {
		s.logger.Errorf(ctx, "Failed to start http server: %v ", err)
		return err
	}

	return nil
}

type Router struct {
	Root *fiber.App
}
