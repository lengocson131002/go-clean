package http

import "github.com/lengocson131002/go-clean/pkg/common"

// DataResponse[T]
// @Description Base Generic Response Body
type DataResponse[T any] struct {
	Status  int    `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func SuccessResponse[T any](data *T) *DataResponse[T] {
	return &DataResponse[T]{
		Status:  common.Success.Status,
		Code:    common.Success.Code,
		Message: common.Success.Message,
		Data:    *data,
	}
}
