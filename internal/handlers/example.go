package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"

	"github.com/evermos/boilerplate-go/internal/domain/example"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/response"
)

// ExampleHandler is the HTTP handler for Example Domain.
type ExampleHandler struct {
	SomeService example.SomeService `inject:"example.someService"`
}

// Router sets up the router for this domain.
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
// @Success 200 {object} response.Base{data=example.SomeEntityFormat}
// @Failure 400 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/example/{id} [get]
func (h *ExampleHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	entity, err := h.SomeService.ResolveByID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, entity)
}
