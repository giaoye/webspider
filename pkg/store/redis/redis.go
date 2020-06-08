package redis

import (
  "github.com/go-redis/redis"
  "sync"
)

var (
  once sync.Once
  client *redis.Client
  
  Nil = redis.Nil
)
func NewClient() *redis.Client {
  return redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
  })
}

func GetClient() *redis.Client {
  once.Do(func() {
    client = NewClient()
  })
  return client
}