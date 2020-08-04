package main

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"

	_ "github.com/go-sql-driver/mysql"
)

const (
	serviceName    = "Evermos/ExampleService"
	serviceVersion = "0.0.1"
)

var (
	config configs.Config
	dbConn infras.MysqlConn
)

func initDb() {
	dbConn.Write = infras.WriteMysqlDB(config)
	dbConn.Read = infras.ReadMysqlDB(config)
}

func initRepositories() {

}

func initServices() {

}

func serveHTTP() {

}
