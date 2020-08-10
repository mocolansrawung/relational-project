package handlers

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/src/services"
)

type ExampleHandler struct {
	ExampleService services.ExampleContract `inject:"service.example"`
}

func (h *ExampleHandler) Router(r chi.Router) {
	r.Get("/health-check", h.Example)
}

func (h *ExampleHandler) Example(w http.ResponseWriter, r *http.Request) {
	status, err := h.ExampleService.Get()
	if err != nil {
		shared.Failed(w, r, shared.WriteFailed(400, "error get status", err.Error()))
		return
	}

	shared.Success(w, r, shared.WriteSuccess(200, "success", status))
}
