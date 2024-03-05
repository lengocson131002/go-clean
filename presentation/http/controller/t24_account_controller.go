package controller

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/pipeline"
	"github.com/lengocson131002/go-clean/pkg/transport/broker"
	"github.com/lengocson131002/go-clean/pkg/transport/http"
)

type T24AccountController struct {
	Logger logger.Logger
	Broker broker.Broker
}

func NewT24AccountController(logger logger.Logger, broker broker.Broker) *T24AccountController {
	return &T24AccountController{
		Logger: logger,
		Broker: broker,
	}
}

func (c *T24AccountController) OpenAccount(ctx *fiber.Ctx) error {
	request := new(domain.OpenAccountRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Logger.Warnf(ctx.UserContext(), "Failed to parse request body : %s", err)
		return fiber.ErrBadRequest
	}

	response, err := pipeline.Send[*domain.OpenAccountRequest, *domain.OpenAccountResponse](ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf(ctx.UserContext(), "Failed to open account : %s", err)
		return err
	}

	httpResp := http.SuccessResponse[*domain.OpenAccountResponse](response)
	return ctx.Status(httpResp.Result.Status).JSON(httpResp)
}

func (c *T24AccountController) OpenAccountRpc(ctx *fiber.Ctx) error {
	request := new(domain.OpenAccountRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Logger.Warnf(ctx.UserContext(), "Failed to parse request body : %s", err)
		return fiber.ErrBadRequest
	}

	reqJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	msg, err := c.Broker.PublishAndReceive(
		"go.test.clean.request",
		&broker.Message{Body: reqJson},
		broker.WithPublishReplyToTopic("go.test.clean.reply"),
		broker.WithReplyConsumerGroup("go.test.reply"),
	)

	if err != nil {
		return err
	}

	var kMsg broker.Response[*domain.OpenAccountResponse]
	err = json.Unmarshal(msg.Body, &kMsg)
	if err != nil {
		return err
	}

	var httpResp http.Response[*domain.OpenAccountResponse]
	if kMsg.Result.Code != "0" {
		httpResp = http.Response[*domain.OpenAccountResponse]{
			Result: http.Result{
				Status:  kMsg.Result.Status,
				Code:    kMsg.Result.Code,
				Message: kMsg.Result.Message,
			},
			Data: nil,
		}
	} else {
		httpResp = http.SuccessResponse[*domain.OpenAccountResponse](kMsg.Data)
	}

	return ctx.Status(httpResp.Result.Status).JSON(httpResp)

}
