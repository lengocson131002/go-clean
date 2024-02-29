package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	dErrors "github.com/lengocson131002/go-clean/pkg/errors"
	"github.com/lengocson131002/go-clean/pkg/transport/http"
)

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	// default message
	msg := http.DefaultErrorResponse
	msg.Result.Message = err.Error()

	// retrieve the custom status code if it's an fiber.*Error
	var e *fiber.Error
	if errors.As(err, &e) {
		msg.Result = http.Result{
			Status:  e.Code,
			Code:    fmt.Sprintf("%v", e.Code),
			Message: e.Message,
		}
	}

	var businessErr *dErrors.DomainError
	if errors.As(err, &businessErr) {
		msg.Result = http.Result{
			Status:  businessErr.Status,
			Code:    businessErr.Code,
			Message: businessErr.Message,
		}
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return ctx.Status(msg.Result.Status).JSON(msg)
}
