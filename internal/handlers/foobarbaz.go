package handlers

import (
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/foobarbaz"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

// FooBarBazHandler is the HTTP handler for FooBarBaz domain.
type FooBarBazHandler struct {
	FooService foobarbaz.FooService
}

// ProvideFooBarBazHandler is the provider for this handler.
func ProvideFooBarBazHandler(fooService foobarbaz.FooService) FooBarBazHandler {
	return FooBarBazHandler{
		FooService: fooService,
	}
}

// Router sets up the router for this domain.
func (h *FooBarBazHandler) Router(r chi.Router) {
	r.Route("/foobarbaz", func(r chi.Router) {
		r.Get("/foo/{id}", h.ResolveFooByID)
	})
}

// ResolveFooByID resolves a Foo by its ID.
// @Summary Resolve Foo by ID
// @Description This endpoint resolves a Foo by its ID.
// @Tags foobarbaz: foo
// @Param id path string true "The Foo's identifier."
// @Produce json
// @Success 200 {object} response.Base{data=foobarbaz.FooFormat}
// @Failure 400 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/foobarbaz/foo/{id} [get]
func (h *FooBarBazHandler) ResolveFooByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	foo, err := h.FooService.ResolveByID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, foo)
}
