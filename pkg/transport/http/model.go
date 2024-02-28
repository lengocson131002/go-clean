package http

// DataResponse[T]
// @Description Base Generic Response Body
type DataResponse[T any] struct {
	Status  int    `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func SuccessResponse[T any](data T) DataResponse[T] {
	var sucRes = DefaultSuccessResponse
	return DataResponse[T]{
		Status:  sucRes.Status,
		Code:    sucRes.Code,
		Message: sucRes.Message,
		Data:    data,
	}
}
