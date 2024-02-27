package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	dErrors "github.com/lengocson131002/go-clean/pkg/errors"
	"github.com/lengocson131002/go-clean/pkg/http"
)

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	// default message
	msg := http.DefaultErrorResponse
	msg.Message = err.Error()

	// retrieve the custom status code if it's an fiber.*Error
	var e *fiber.Error
	if errors.As(err, &e) {
		msg = http.DataResponse[interface{}]{
			Status:  e.Code,
			Code:    e.Code,
			Message: e.Message,
		}
	}

	var businessErr *dErrors.DomainError
	if errors.As(err, &businessErr) {
		msg = http.DataResponse[interface{}]{
			Status:  businessErr.Status,
			Code:    businessErr.Code,
			Message: businessErr.Message,
		}
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return ctx.Status(msg.Status).JSON(msg)
}
