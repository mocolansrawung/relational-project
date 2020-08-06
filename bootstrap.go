package main

import (
	"log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/container"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/src/handlers"
	"github.com/evermos/boilerplate-go/src/repositories"
	"github.com/evermos/boilerplate-go/src/services"
	_ "github.com/go-sql-driver/mysql"
)

const (
	serviceName    = "Evermos/ExampleService"
	serviceVersion = "0.0.1"
)

func init() {
	container := container.NewContainer()
	config := configs.Get()
	db := infras.MysqlConn{Write: infras.WriteMysqlDB(*config), Read: infras.ReadMysqlDB(*config)}
	container.Register("config", config)
	container.Register("db", &db)
	container.Register("repo", new(repositories.Repository))
	container.Register("service", new(services.Service))

	h := handlers.Handler{}
	container.Register("handler", &h)

	err := container.Start()
	if err != nil {
		log.Fatalf("error starting container : %v", err)
	}

	h.TestHandler()
}
