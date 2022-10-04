package services

import (
	"os"
	"fmt"
	"log"
	"github.com/joho/godotenv"
	redis "github.com/go-redis/redis/v7"
)

var  client *redis.Client

func InitRedis() {
	//Initializing redis
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

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
     panic(err)
  }
}

func GetRedis() *redis.Client {
	return client
}