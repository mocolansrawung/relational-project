package main

import (
	"log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/container"
	"github.com/evermos/boilerplate-go/events/example"
	"github.com/evermos/boilerplate-go/infras"
	routers "github.com/evermos/boilerplate-go/router"
	"github.com/evermos/boilerplate-go/src/handlers"
	"github.com/evermos/boilerplate-go/src/repositories"
	"github.com/evermos/boilerplate-go/src/services"

	_ "github.com/go-sql-driver/mysql"
)

func registry() *container.ServiceRegistry {
	c := container.NewContainer()
	config = configs.Get()
	db = &infras.MysqlConn{Write: infras.WriteMysqlDB(*config), Read: infras.ReadMysqlDB(*config)}
	c.Register("config", config)
	c.Register("db", db)

	// Repository
	c.Register("repository.example", new(repositories.ExampleRepository))

	// Service
	c.Register("service.example", new(services.ExampleService))

	// Handler
	c.Register("handler.example", new(handlers.ExampleHandler))

	// Event
	c.Register("consumer.example", new(example.EventConsumer))

	router = &routers.Router{}
	c.Register("router", router)
	err := c.Start()
	if err != nil {
		log.Fatalln(err)
	}

	return c
}
