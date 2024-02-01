package http

import "github.com/lengocson131002/go-clean/pkg/common"

var (
	DefaultSuccessResponse = DataResponse[interface{}]{
		Status:  common.Success.Status,
		Code:    common.Success.Code,
		Message: common.Success.Message,
		Data:    nil,
	}

	DefaultErrorResponse = DataResponse[interface{}]{
		Status:  common.ErrInternalServer.Status,
		Code:    common.ErrInternalServer.Status,
		Message: common.ErrInternalServer.Message,
		Data:    nil,
	}
)
