package model

import (
	"github.com/gofiber/fiber/v2"
	response "github.com/lengocson131002/go-clean/internal/domain/response"
)

type DataResponse[T any] struct {
	Status  int    `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func (g *DataResponse[T]) JSON(c *fiber.Ctx) error {
	return c.Status(g.Status).JSON(g)
}

func SuccessResponse[T any](data *T) *DataResponse[T] {
	return &DataResponse[T]{
		Status:  200,
		Code:    response.Sucecss.Code,
		Message: response.Sucecss.Message,
		Data:    *data,
	}
}
