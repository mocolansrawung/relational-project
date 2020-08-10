package shared

import (
	"net/http"

	"github.com/go-chi/render"
)

type ResponseFailed struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

type ResponseSuccess struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WriteFailed(code int, message string, err interface{}) ResponseFailed {
	return ResponseFailed{
		Code:    code,
		Message: message,
		Error:   err,
	}
}

func WriteSuccess(code int, message string, data interface{}) ResponseSuccess {
	return ResponseSuccess{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func Success(w http.ResponseWriter, r *http.Request, resp ResponseSuccess) {
	render.Status(r, resp.Code)
	render.JSON(w, r, resp)
}

func Failed(w http.ResponseWriter, r *http.Request, resp ResponseFailed) {
	render.Status(r, resp.Code)
	render.JSON(w, r, resp)
}
