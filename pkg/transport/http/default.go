package http

var (
	DefaultSuccessResponse = DataResponse[interface{}]{
		Status:  200,
		Code:    0,
		Message: "Success",
		Data:    nil,
	}

	DefaultErrorResponse = DataResponse[interface{}]{
		Status:  500,
		Code:    500,
		Message: "Internal Server Error",
		Data:    nil,
	}
)
