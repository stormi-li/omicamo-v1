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

	cacheGetCallback := func(key string) string {
		res, _ := c.RedisClient.Get(context.Background(), key).Result()
		return res
	}
	databaseGetCallback := func(key string) string { return "" }
	cacheSetCallback := func(key, value string) {
		c.RedisClient.Set(context.Background(), key, value, 5*time.Minute)
	}
	databaseSetCallback := func(key, value string) {}

	cache.AddCallback(cacheGetCallback, databaseGetCallback, cacheSetCallback, databaseSetCallback)

	key := "name"
	value := "stormi-li"
	
	res := cache.Get(key)
	if res == omicamo.NullString {
		cache.Set(key, value)
	}
	res = cache.Get(key)
	fmt.Println(res)
}
