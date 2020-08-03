package main

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	serviceName    = "Evermos/ExampleService"
	serviceVersion = "0.0.1"
)

var (
	config      configs.Config
	writeDbConn *sqlx.DB
	readDbConn  *sqlx.DB
)

func initDb() {
	writeDbConn = infras.WriteMysqlDB(config)
	readDbConn = infras.ReadMysqlDB(config)
}

func initRepositories() {

}

func initServices() {

}

func serveHTTP() {

}
