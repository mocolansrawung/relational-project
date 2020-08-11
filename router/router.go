package router

import (
	"github.com/evermos/boilerplate-go/src/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Router struct {
	ExampleHandler *handlers.ExampleHandler `inject:"handler.example"`
}

func (r *Router) Route() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	r.ExampleHandler.Router(mux)

	return mux
}
