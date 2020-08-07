package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/src/services"
)

type Handler struct {
	ExampleService services.ExampleContract `inject:"service.example"`
}

func (h *Handler) Router(r chi.Router) {
	r.Get("/health-check", h.Example)
}

func (h *Handler) Example(w http.ResponseWriter, r *http.Request) {
	status, err := h.ExampleService.Get()
	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, shared.NewErrorResponse(err, "error get status", 0))
		return
	}

	render.JSON(w, r, shared.NewResponse(status, "success", 0))
}
