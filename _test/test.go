package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omicamo-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	c := omicamo.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	cache := c.NewCache()
	cache.SetCacheCallback(func(key string, redisClient *redis.Client) string {
		res, _ := redisClient.Get(context.Background(), key).Result()
		return res
	}, func(key, value string, redisClient *redis.Client) {
		if value == omicamo.NullString {
			redisClient.Set(context.Background(), key, value, 1*time.Minute)
			return
		}
		redisClient.Set(context.Background(), key, value, 5*time.Minute)
	})
	cache.SetDatabaseCallback(func(key string) string {
		return key
	}, func(key, value string) {})

	cache.Set("key1", "fsfsfsf")
	res := cache.Get("key1")
	fmt.Println(res)
}
