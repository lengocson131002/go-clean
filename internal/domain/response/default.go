package response

type DomainResponse struct {
	Status  int    // http status mapping
	Code    int    // domain response code
	Message string // domain response message
}

func (err *DomainResponse) Error() string {
	return err.Message
}

var (
	// COMMON ERROR
	Sucecss = &DomainResponse{
		Status:  200,
		Code:    0,
		Message: "Success",
	}

	ErrorCancelled = &DomainResponse{
		Status:  499,
		Code:    1,
		Message: "Cancelled",
	}

	ErrorUnknown = &DomainResponse{
		Status:  500,
		Code:    2,
		Message: "Unknow",
	}

	ErrBadRequest = &DomainResponse{
		Status:  400,
		Code:    3,
		Message: "Bad request",
	}

	ErrorDeadlineExceeded = &DomainResponse{
		Status:  504,
		Code:    4,
		Message: "Deadline exceeded",
	}

	ErrNotFound = &DomainResponse{
		Status:  404,
		Code:    5,
		Message: "Not found",
	}

	ErrAlreadyExists = &DomainResponse{
		Status:  409,
		Code:    6,
		Message: "Already exists",
	}

	ErrPermissionDenied = &DomainResponse{
		Status:  403,
		Code:    7,
		Message: "Permission denied",
	}

	ErrorResourceExaushted = &DomainResponse{
		Status:  429,
		Code:    8,
		Message: "Too many requests",
	}

	ErrorFailedPrecondition = &DomainResponse{
		Status:  400,
		Code:    9,
		Message: "Failed precondition",
	}

	ErrorAborted = &DomainResponse{
		Status:  409,
		Code:    10,
		Message: "Aborded",
	}

	ErrorOutOfRange = &DomainResponse{
		Status:  400,
		Code:    11,
		Message: "Out of range",
	}

	ErrorUnimplemented = &DomainResponse{
		Status:  501,
		Code:    12,
		Message: "Unimlemented",
	}

	ErrInternalServer = &DomainResponse{
		Status:  500,
		Code:    13,
		Message: "Internal server error",
	}

	ErrUnavailable = &DomainResponse{
		Status:  503,
		Code:    14,
		Message: "Internal server error",
	}

	ErrorDataLoss = &DomainResponse{
		Status:  500,
		Code:    15,
		Message: "Internal server error",
	}

	ErrUnauthorized = &DomainResponse{
		Code:    16,
		Message: "Unauthorized",
	}

	// DOMAIN CUSTOM ERROR
	ErrorAccountNotFound = &DomainResponse{
		Status:  400,
		Code:    100,
		Message: "User not found",
	}

	ErrorAccountExisted = &DomainResponse{
		Status:  400,
		Code:    101,
		Message: "User ID already existed",
	}
)
