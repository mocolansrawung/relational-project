package router

import (
	"github.com/evermos/boilerplate-go/internals/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Router struct {
	ExampleHandler *handlers.ExampleHandler `inject:"handler.example"`
}

func (r *Router) NewRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(rc chi.Router) {
		r.ExampleHandler.Router(rc)
	})

	return mux
}
