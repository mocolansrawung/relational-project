package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"

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

var config *configs.Config

type Router struct {
	ExampleHandler *handlers.ExampleHandler `inject:"handler.example"`
}

func registry() *container.ServiceRegistry {
	c := container.NewContainer()
	config = configs.Get()
	db := infras.MysqlConn{Write: infras.WriteMysqlDB(*config), Read: infras.ReadMysqlDB(*config)}
	c.Register("config", config)
	c.Register("db", &db)

	// Repository
	c.Register("repository.example", new(repositories.ExampleRepository))

	// Service
	c.Register("service.example", new(services.ExampleService))

	// Handler
	c.Register("handler.example", new(handlers.ExampleHandler))

	return c
}

func ServeHTTP() error {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	c := registry()

	router := Router{}
	c.Register("router", &router)

	if err := c.Start(); err != nil {
		log.Fatalln(err)
	}
	router.ExampleHandler.Router(mux)

	return http.ListenAndServe(":"+config.Port, mux)
}
