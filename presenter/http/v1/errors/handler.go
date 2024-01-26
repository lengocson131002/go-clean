package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/internal/domain/response"
	"github.com/lengocson131002/go-clean/presenter/http/v1/model"
)

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	msg := model.DefaultErrorResponse
	msg.Message = err.Error()

	// trieve the custom status code if it's an fiber.*Error
	var e *fiber.Error
	if errors.As(err, &e) {
		msg = model.DataResponse[interface{}]{
			Status:  e.Code,
			Code:    e.Code,
			Message: e.Message,
		}
	}
	var customErr *response.DomainResponse
	if errors.As(err, &customErr) {
		msg = model.DataResponse[interface{}]{
			Status:  customErr.Status,
			Code:    customErr.Code,
			Message: customErr.Message,
		}
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return msg.JSON(ctx)
}
