package main

//go:generate go run github.com/swaggo/swag/cmd/swag init

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/di"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/internal/domain/example"
	"github.com/evermos/boilerplate-go/internal/handlers"
	"github.com/evermos/boilerplate-go/router"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/evermos/boilerplate-go/transport/http"
	"github.com/rs/zerolog/log"
)

var (
	db         *infras.MySQLConn
	config     *configs.Config
	httpRouter *router.Router
)

func main() {
	// Initialize logger
	logger.InitLogger()

	// Initialize config
	config = configs.Get()

	// Set desired log level
	logger.SetLogLevel(config)

	// Initialize container
	container := di.NewContainer()
	db = &infras.MySQLConn{
		Write: infras.CreateMySQLWriteConn(*config),
		Read:  infras.CreateMySQLReadConn(*config)}
	container.Register("config", config)
	container.Register("db", db)

	// Domain - Example
	container.Register("example.someRepository", new(example.SomeRepositoryMySQL))
	container.Register("example.someService", new(example.SomeServiceImpl))

	// Handlers
	container.Register("handler.example", new(handlers.ExampleHandler))

	httpRouter = &router.Router{}
	container.Register("router", httpRouter)
	err := container.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	// Setup HTTP server
	http := http.Http{
		DB:     db,
		Config: config,
		Router: httpRouter.NewRouter()}
	http.GracefulShutdown(*container)

	// Run server
	http.Serve()
}
