package main

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/delivery/http"
	"github.com/evermos/boilerplate-go/infras"
	routers "github.com/evermos/boilerplate-go/router"
)

var (
	db     *infras.MysqlConn
	config *configs.Config
	router *routers.Router
)

func main() {
	registry()
	http := http.Http{DB: db, Config: config, Router: router.NewRouter()}
	http.GracefulShutdown()
	http.Serve()
}
