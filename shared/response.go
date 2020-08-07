package shared

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(data interface{}, message string, code int) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewErrorResponse(err error, message string, code int) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Error:   err.Error(),
	}
}
