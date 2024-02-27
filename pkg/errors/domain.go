package errors

import "net/http"

type DomainError struct {
	Status  int    // http status mapping
	Code    int    // domain error code
	Message string // domain error message
}

func (err *DomainError) Error() string {
	return err.Message
}

var (
	DomainValidationError = &DomainError{
		Status:  http.StatusBadRequest,
		Code:    1,
		Message: "Validation error",
	}
)
