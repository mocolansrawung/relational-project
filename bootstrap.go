package main

import (
	"log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/container"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/src/handlers"
	"github.com/evermos/boilerplate-go/src/repositories"
	"github.com/evermos/boilerplate-go/src/services"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

const (
	serviceName    = "Evermos/ExampleService"
	serviceVersion = "0.0.1"
)

type Router struct {
	Handler *handlers.Handler `inject:"handler"`
}

func registry() *container.ServiceRegistry {
	c := container.NewContainer()
	config := configs.Get()
	db := infras.MysqlConn{Write: infras.WriteMysqlDB(*config), Read: infras.ReadMysqlDB(*config)}
	c.Register("config", config)
	c.Register("db", &db)

	// Repository
	c.Register("repository.example", new(repositories.ExampleRepository))

	// Service
	c.Register("service.example", new(services.ExampleService))

	//Handler
	c.Register("handler", new(handlers.Handler))

	return c
}

func ServeHTTP() *chi.Mux {
	mux := chi.NewRouter()
	c := registry()

	router := Router{}
	c.Register("router", &router)

	if err := c.Start(); err != nil {
		log.Fatalln(err)
	}
	router.Handler.Router(mux)

	return mux
}
