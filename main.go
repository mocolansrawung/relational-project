package main

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	routers "github.com/evermos/boilerplate-go/router"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/evermos/boilerplate-go/transport/http"
)

var (
	db     *infras.MysqlConn
	config *configs.Config
	router *routers.Router
)

func main() {
	logger.InitLogger()
	container := registry()
	http := http.Http{DB: db, Config: config, Router: router.NewRouter()}
	http.GracefulShutdown(*container)
	http.Serve()
}
