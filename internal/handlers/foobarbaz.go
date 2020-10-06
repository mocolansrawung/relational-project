package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
		r.Post("/foo", h.CreateFoo)
		r.Get("/foo/{id}", h.ResolveFooByID)
	})
}

// CreateFoo creates a new Foo.
// @Summary Create a new Foo.
// @Description This endpoint creates a new Foo.
// @Tags foobarbaz/foo
// @Param foo body foobarbaz.FooRequestFormat true "The Foo to be created."
// @Produce json
// @Success 201 {object} response.Base{data=foobarbaz.FooResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/foobarbaz/foo [post]
func (h *FooBarBazHandler) CreateFoo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat foobarbaz.FooRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, _ := uuid.NewV4() // TODO: read from context

	foo, err := h.FooService.Create(requestFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, foo)
}

// ResolveFooByID resolves a Foo by its ID.
// @Summary Resolve Foo by ID
// @Description This endpoint resolves a Foo by its ID.
// @Tags foobarbaz/foo
// @Param id path string true "The Foo's identifier."
// @Param withItems query string false "Fetch with items, default false."
// @Produce json
// @Success 200 {object} response.Base{data=foobarbaz.FooResponseFormat}
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

	withItems, _ := strconv.ParseBool(r.URL.Query().Get("withItems"))

	foo, err := h.FooService.ResolveByID(id, withItems)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, foo)
}
