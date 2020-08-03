package configs

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseHost     string `envconfig:"database_host"`
	DatabaseUsername string `envconfig:"database_username"`
	DatabasePassword string `envconfig:"database_password"`
	DatabaseName     string `envconfig:"database_name"`
}

var conf Config

func Get() Config {
	once := sync.Once{}
	once.Do(func() {
		err := envconfig.Process("", &conf)
		if err != nil {
			log.Fatalf("err : %v", err)
		}
	})

	return conf
}
