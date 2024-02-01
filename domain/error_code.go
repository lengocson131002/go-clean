package domain

import "github.com/lengocson131002/go-clean/pkg/common"

var (
	// DOMAIN CUSTOM ERROR
	ErrorAccountNotFound = &common.InternalResponse{
		Status:  400,
		Code:    100,
		Message: "User not found",
	}

	ErrorAccountExisted = &common.InternalResponse{
		Status:  400,
		Code:    101,
		Message: "User ID already existed",
	}
)
