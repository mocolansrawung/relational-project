package infras

import (
	"fmt"
	"log"

	"github.com/evermos/boilerplate-go/configs"

	"github.com/jmoiron/sqlx"
)

const (
	maxIdleConnection = 10
	maxOpenConnection = 10
)

// WriteMysqlDB - function for creating database connection for write-access
func WriteMysqlDB(config configs.Config) *sqlx.DB {
	return CreateDBConnection(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=%s&parseTime=true",
		config.DatabaseUsername, config.DatabasePassword, config.DatabaseHost, config.DatabaseName, config.DatabaseTimeZone))

}

// ReadMysqlDB function for creating database connection for read-access
func ReadMysqlDB(config configs.Config) *sqlx.DB {
	return CreateDBConnection(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=%s&parseTime=true",
		config.DatabaseUsername, config.DatabasePassword, config.DatabaseHost, config.DatabaseName, config.DatabaseTimeZone))

}

// CreateDBConnection function for creating database connection
func CreateDBConnection(descriptor string) *sqlx.DB {
	db, err := sqlx.Connect("mysql", descriptor)
	if err != nil {
		log.Fatalf("error connecting to DB: %s", descriptor)
	}
	db.SetMaxIdleConns(maxIdleConnection)
	db.SetMaxOpenConns(maxOpenConnection)

	return db
}
