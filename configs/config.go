package configs

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	//Mysql
	DatabaseHost     string `envconfig:"database_host"`
	DatabaseUsername string `envconfig:"database_username"`
	DatabasePassword string `envconfig:"database_password"`
	DatabaseName     string `envconfig:"database_name"`
	DatabaseTimeZone string `envconfig:"database_time_zone"`

	//Redis
	RedisHost     string `envconfig:"redis_host"`
	RedisPort     string `envconfig:"redis_port"`
	RedisPassword string `envconfig:"redis_password"`
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
