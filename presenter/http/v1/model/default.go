package model

import response "github.com/lengocson131002/go-clean/internal/domain/response"

var (
	DefaultSuccessResponse = DataResponse[interface{}]{
		Status:  200,
		Code:    response.Sucecss.Code,
		Message: response.Sucecss.Message,
		Data:    nil,
	}

	DefaultErrorResponse = DataResponse[interface{}]{
		Status:  500,
		Code:    response.ErrInternalServer.Code,
		Message: response.ErrInternalServer.Message,
		Data:    nil,
	}
)
