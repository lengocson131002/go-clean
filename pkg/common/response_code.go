package common

type InternalResponse struct {
	Status  int    // http status mapping
	Code    int    // domain error code
	Message string // domain error message
}

func (err *InternalResponse) Error() string {
	return err.Message
}

var (
	// COMMON ERROR
	Success = &InternalResponse{
		Status:  200,
		Code:    0,
		Message: "Success",
	}

	ErrorCancelled = &InternalResponse{
		Status:  499,
		Code:    1,
		Message: "Cancelled",
	}

	ErrorUnknown = &InternalResponse{
		Status:  500,
		Code:    2,
		Message: "Unknown",
	}

	ErrBadRequest = &InternalResponse{
		Status:  400,
		Code:    3,
		Message: "Bad request",
	}

	ErrorDeadlineExceeded = &InternalResponse{
		Status:  504,
		Code:    4,
		Message: "Deadline exceeded",
	}

	ErrNotFound = &InternalResponse{
		Status:  404,
		Code:    5,
		Message: "Not found",
	}

	ErrAlreadyExists = &InternalResponse{
		Status:  409,
		Code:    6,
		Message: "Already exists",
	}

	ErrPermissionDenied = &InternalResponse{
		Status:  403,
		Code:    7,
		Message: "Permission denied",
	}

	ErrorResourceExaushted = &InternalResponse{
		Status:  429,
		Code:    8,
		Message: "Too many requests",
	}

	ErrorFailedPrecondition = &InternalResponse{
		Status:  400,
		Code:    9,
		Message: "Failed precondition",
	}

	ErrorAborted = &InternalResponse{
		Status:  409,
		Code:    10,
		Message: "Aborded",
	}

	ErrorOutOfRange = &InternalResponse{
		Status:  400,
		Code:    11,
		Message: "Out of range",
	}

	ErrorUnimplemented = &InternalResponse{
		Status:  501,
		Code:    12,
		Message: "Unimlemented",
	}

	ErrInternalServer = &InternalResponse{
		Status:  500,
		Code:    13,
		Message: "Internal server error",
	}

	ErrUnavailable = &InternalResponse{
		Status:  503,
		Code:    14,
		Message: "Internal server error",
	}

	ErrorDataLoss = &InternalResponse{
		Status:  500,
		Code:    15,
		Message: "Internal server error",
	}

	ErrUnauthorized = &InternalResponse{
		Code:    16,
		Message: "Unauthorized",
	}
)
