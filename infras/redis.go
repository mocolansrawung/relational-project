package infras

import (
	"fmt"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/go-redis/redis"
)

//RedisNewClient create new instance of redis
func RedisNewClient(config configs.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong, err)

	return client
}
