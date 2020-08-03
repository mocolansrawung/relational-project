package configs

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	MaxIdleConnection = 10
	MaxOpenConnection = 10
)

// WriteMysqlDB - function for creating database connection for write-access
func WriteMysqlDB() *sqlx.DB {
	return CreateDBConnection(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=%s&parseTime=true",
		os.Getenv("WRITE_DB_USER"), os.Getenv("WRITE_DB_PASSWORD"), os.Getenv("WRITE_DB_HOST"), os.Getenv("WRITE_DB_NAME"), url.QueryEscape("Asia/Jakarta")))

}

// ReadMysqlDB function for creating database connection for read-access
func ReadMysqlDB() *sqlx.DB {
	return CreateDBConnection(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=%s&parseTime=true",
		os.Getenv("READ_DB_USER"), os.Getenv("READ_DB_PASSWORD"), os.Getenv("READ_DB_HOST"), os.Getenv("READ_DB_NAME"), url.QueryEscape("Asia/Jakarta")))

}

// CreateDBConnection function for creating database connection
func CreateDBConnection(descriptor string) *sqlx.DB {
	db, err := sqlx.Connect("mysql", descriptor)
	if err != nil {
		log.Fatalf("error connecting to DB: %s", descriptor)
	}
	db.SetMaxIdleConns(MaxIdleConnection)
	db.SetMaxOpenConns(MaxOpenConnection)

	return db
}