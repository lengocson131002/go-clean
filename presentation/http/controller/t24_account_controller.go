package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/http"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/pipeline"
)

type T24AccountController struct {
	Logger logger.Logger
}

func NewT24AccountController(logger logger.Logger) *T24AccountController {
	return &T24AccountController{
		Logger: logger,
	}
}

func (c *T24AccountController) OpenAccount(ctx *fiber.Ctx) error {
	request := new(domain.OpenAccountRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Logger.Warn(ctx.UserContext(), "Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := pipeline.Send[*domain.OpenAccountRequest, *domain.OpenAccountResponse](ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warn(ctx.UserContext(), "Failed to open account : %+v", err)
		return err
	}

	httpResp := http.SuccessResponse[*domain.OpenAccountResponse](response)
	return ctx.Status(httpResp.Status).JSON(httpResp)
}
