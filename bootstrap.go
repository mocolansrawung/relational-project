package main

import (
	"log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/di"
	"github.com/evermos/boilerplate-go/events/example"
	"github.com/evermos/boilerplate-go/events/producer"
	"github.com/evermos/boilerplate-go/infras"
	exampleDomain "github.com/evermos/boilerplate-go/internal/domain/example"
	"github.com/evermos/boilerplate-go/internal/handlers"
	routers "github.com/evermos/boilerplate-go/router"

	_ "github.com/go-sql-driver/mysql"
)

func registry() *di.ServiceRegistry {
	c := di.NewContainer()
	config = configs.Get()
	db = &infras.MysqlConn{Write: infras.WriteMysqlDB(*config), Read: infras.ReadMysqlDB(*config)}
	c.Register("config", config)
	c.Register("db", db)

	// Domain - Example
	c.Register("example.someRepository", new(exampleDomain.SomeRepositoryMySQL))
	c.Register("example.someService", new(exampleDomain.SomeServiceImpl))

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
