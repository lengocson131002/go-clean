package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/pkg/common"
	"github.com/lengocson131002/go-clean/pkg/http"
)

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	msg := http.DefaultErrorResponse
	msg.Message = err.Error()

	// trieve the custom status code if it's an fiber.*Error
	var e *fiber.Error
	if errors.As(err, &e) {
		msg = http.DataResponse[interface{}]{
			Status:  e.Code,
			Code:    e.Code,
			Message: e.Message,
		}
	}
	var customErr *common.InternalResponse
	if errors.As(err, &customErr) {
		msg = http.DataResponse[interface{}]{
			Status:  customErr.Status,
			Code:    customErr.Code,
			Message: customErr.Message,
		}
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return ctx.Status(msg.Status).JSON(msg)
}
