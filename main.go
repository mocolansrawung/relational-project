package main

//go:generate go run github.com/swaggo/swag/cmd/swag init
//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/evermos/boilerplate-go/transport/http"
)

var (
	db     *infras.MySQLConn
	config *configs.Config
)

func main() {
	// Initialize logger
	logger.InitLogger()

	// Initialize config
	config = configs.Get()

	// Set desired log level
	logger.SetLogLevel(config)

	// Wire everything up
	httpRouter := InitializeService()

	// Setup HTTP server
	http := http.HTTP{
		DB:     db,
		Config: config,
		Router: httpRouter}

	// Run server
	http.SetupAndServe()
}
