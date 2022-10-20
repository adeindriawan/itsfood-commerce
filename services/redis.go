package services

import (
	"os"
	"fmt"
	redis "github.com/go-redis/redis/v7"
)

var  client *redis.Client

func InitRedis() {
  dsn := os.Getenv("REDIS_HOST")
  if len(dsn) == 0 {
     dsn = "localhost:6379"
  }
	fmt.Println(dsn)
  client = redis.NewClient(&redis.Options{
     Addr: dsn, //redis port
  })
  _, errRedis := client.Ping().Result()
  if errRedis != nil {
     panic(errRedis)
  }
}

func GetRedis() *redis.Client {
	return client
}