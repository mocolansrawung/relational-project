package main

import (
	"log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/di"
	"github.com/evermos/boilerplate-go/events/example"
	"github.com/evermos/boilerplate-go/events/producer"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/internal/handlers"
	"github.com/evermos/boilerplate-go/internal/repositories"
	"github.com/evermos/boilerplate-go/internal/services"
	routers "github.com/evermos/boilerplate-go/router"

	_ "github.com/go-sql-driver/mysql"
)

func registry() *di.ServiceRegistry {
	c := di.NewContainer()
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

	// Producer
	c.Register("producer", new(producer.Producer))

	router = &routers.Router{}
	c.Register("router", router)
	err := c.Start()
	if err != nil {
		log.Fatalln(err)
	}

	return c
}
