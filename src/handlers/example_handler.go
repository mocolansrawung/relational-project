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

// Example :nodoc:
// @Summary Example
// @Description Example
// @Tags Example
// @Produce json
// @Accept json
// @Success 200 {object} shared.Response
// @Router /health-check [get]
func (h *ExampleHandler) Example(w http.ResponseWriter, r *http.Request) {
	status, err := h.ExampleService.Get()
	if err != nil {
		shared.JsonResponse(w, r, shared.NewResponse(400, "error get status", nil, err.Error()))
		return
	}

	shared.JsonResponse(w, r, shared.NewResponse(200, "success", status))
}
