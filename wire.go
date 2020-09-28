//+build wireinject

package main

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/internal/domain/foobarbaz"
	"github.com/evermos/boilerplate-go/internal/handlers"
	"github.com/evermos/boilerplate-go/router"
	"github.com/google/wire"
)

// Wiring for configurations.
var configurations = wire.NewSet(
	configs.Get,
)

// Wiring for persistences.
var persistences = wire.NewSet(
	infras.ProvideMySQLConn,
)

// Wiring for domain FooBarBaz.
var domainFooBarBaz = wire.NewSet(
	// FooService interface and implementation
	foobarbaz.ProvideFooServiceImpl,
	wire.Bind(new(foobarbaz.FooService), new(*foobarbaz.FooServiceImpl)),
	// FooRepository interface and implementation
	foobarbaz.ProvideFooRepositoryMySQL,
	wire.Bind(new(foobarbaz.FooRepository), new(*foobarbaz.FooRepositoryMySQL)),
)

// Wiring for all domains.
var domains = wire.NewSet(
	domainFooBarBaz,
)

// Wiring for HTTP routing.
var routing = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "FooBarBazHandler"),
	handlers.ProvideFooBarBazHandler,
	router.ProvideRouter,
)

// Wiring for everything.
func InitializeService() router.Router {
	wire.Build(
		// configurations
		configurations,
		// persistences
		persistences,
		// domains
		domains,
		// routing
		routing)
	return router.Router{}
}
