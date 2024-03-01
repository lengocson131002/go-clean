package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/pkg/transport/http"
)

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	fRes := http.FailureResponse(err)
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return ctx.Status(fRes.Result.Status).JSON(fRes)
}
