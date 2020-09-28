package router

import (
	"github.com/evermos/boilerplate-go/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type DomainHandlers struct {
	FooBarBazHandler handlers.FooBarBazHandler
}

type Router struct {
	DomainHandlers DomainHandlers
}

func ProvideRouter(domainHandlers DomainHandlers) Router {
	return Router{
		DomainHandlers: domainHandlers,
	}
}

func (r *Router) NewRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(rc chi.Router) {
		r.DomainHandlers.FooBarBazHandler.Router(rc)
	})

	return mux
}
