package configs

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

//Config config struct consist of data that provided from env
type Config struct {
	//Mysql
	WriteDatabaseHost     string `mapstructure:"WRITE_DB_HOST"`
	WriteDatabaseUsername string `mapstructure:"WRITE_DB_USER"`
	WriteDatabasePassword string `mapstructure:"WRITE_DB_PASSWORD"`
	WriteDatabaseName     string `mapstructure:"WRITE_DB_NAME"`
	WriteDatabaseTimeZone string `mapstructure:"WRITE_DB_TIME_ZONE"`

	ReadDatabaseHost     string `mapstructure:"READ_DB_HOST"`
	ReadDatabaseUsername string `mapstructure:"READ_DB_USER"`
	ReadDatabasePassword string `mapstructure:"READ_DB_PASSWORD"`
	ReadDatabaseName     string `mapstructure:"READ_DB_NAME"`
	ReadDatabaseTimeZone string `mapstructure:"READ_DB_TIME_ZONE"`

	//Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
}

var conf Config

//Get are responsible to load env and get data an return the struct
func Get() Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	once := sync.Once{}
	once.Do(func() {
		err = viper.Unmarshal(&conf)
		if err != nil {
			log.Fatalln(err)
		}
	})

	return conf
}
