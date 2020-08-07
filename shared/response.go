package shared

import (
	"net/http"
	"reflect"

	"github.com/go-chi/render"
)

type Response struct {
	Code    int         `json:"-"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewResponse(code int, message string, data interface{}, err ...interface{}) Response {
	r := Response{
		Code:    code,
		Message: message,
	}
	if data != nil || reflect.ValueOf(data).Kind() == reflect.Ptr {
		r.Data = data
	}

	if err != nil || reflect.ValueOf(err).Kind() == reflect.Ptr {
		r.Error = err[0]
	}

	return r
}

func JsonResponse(w http.ResponseWriter, r *http.Request, resp Response) {
	render.Status(r, resp.Code)
	render.JSON(w, r, resp)
}
