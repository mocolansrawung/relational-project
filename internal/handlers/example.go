package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"

	"github.com/evermos/boilerplate-go/internal/domain/example"
	"github.com/evermos/boilerplate-go/shared"
)

type ExampleHandler struct {
	SomeService example.SomeService `inject:"example.someService"`
}

func (h *ExampleHandler) Router(r chi.Router) {
	r.Route("/example", func(r chi.Router) {
		r.Get("/{id}", h.ResolveByID)
	})
}

// ResolveByID is a sample endpoint for resolving an entity by its ID
// @Summary Resolve by ID
// @Description An example endpoint that resolves an entity by its ID
// @Tags example
// @Param id path string true "The entity's identifier."
// @Produce json
// @Success 200 {object} shared.ResponseSuccess
// @Failure 400 {object} shared.ResponseFailed
// @Router /v1/example/{id} [get]
func (h *ExampleHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		shared.Failed(w, r, shared.WriteFailed(400, "error resolving entity", err.Error()))
		return
	}

	entity, err := h.SomeService.ResolveByID(id)
	if err != nil {
		shared.Failed(w, r, shared.WriteFailed(400, "error resolving entity", err.Error()))
		return
	}

	shared.Success(w, r, shared.WriteSuccess(200, "success", entity))
}
