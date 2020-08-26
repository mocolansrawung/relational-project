package handlers

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/evermos/boilerplate-go/internals/services"
	"github.com/evermos/boilerplate-go/shared"
)

type ExampleHandler struct {
	ExampleService services.ExampleContract `inject:"service.example"`
}

func (h *ExampleHandler) Router(r chi.Router) {
	r.Route("/example", func(r chi.Router) {
		r.Get("/", h.Example)
	})
}

// Example godoc
// @Summary Healthcheck
// @Description Healthcheck Endpoint
// @Tags Healthcheck
// @Produce json
// @Accept json
// @Success 200 {object} shared.ResponseSuccess
// @Failure 400 {object} shared.ResponseFailed
// @Router /health-check [get]
func (h *ExampleHandler) Example(w http.ResponseWriter, r *http.Request) {
	status, err := h.ExampleService.Get()
	if err != nil {
		shared.Failed(w, r, shared.WriteFailed(400, "error get status", err.Error()))
		return
	}

	shared.Success(w, r, shared.WriteSuccess(200, "success", status))
}
