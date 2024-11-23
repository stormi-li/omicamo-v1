# Omicamo 缓存安全框架
**作者**: stormi-li  
**Email**: 2785782829@qq.com  
## 简介

**Omicamo** 是一个缓存安全框架，专注于解决分布式环境下的缓存一致性、缓存击穿以及缓存穿透问题。通过集成分布式锁与空值存储机制，**Omicamo** 提供了强大的缓存管理功能。



## 功能

- **缓存一致性**：通过分布式锁保证缓存与数据库的一致性。
- **缓存击穿保护**：避免高并发请求集中访问数据库。
- **缓存穿透防御**：通过存储空值拦截非法请求，减少对数据库的无效访问。
## 教程
### 安装
```shell
go get github.com/stormi-li/omicamo-v1
```
### 使用
```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omicamo-v1"
)

func main() {
	// 创建 Omicamo 客户端，连接 Redis
	c := omicamo.NewClient(&redis.Options{Addr: "localhost:6379"})
	cache := c.NewCache()
	
	// 获取缓存的回调函数
	cacheGetCallback := func(key string) string {
		res, _ := c.RedisClient.Get(context.Background(), key).Result()
		return res
	}

	// 获取数据库的回调函数（当缓存没有命中时使用）
	databaseGetCallback := func(key string) string { return "" }

	// 设置缓存的回调函数
	cacheSetCallback := func(key, value string) {
		c.RedisClient.Set(context.Background(), key, value, 5*time.Minute) // 设置缓存有效期5分钟
	}

	// 设置数据库的回调函数（当缓存修改时写入数据库）
	databaseSetCallback := func(key, value string) {}

	// 注册回调函数
	cache.AddCallback(cacheGetCallback, databaseGetCallback, cacheSetCallback, databaseSetCallback)

	// 缓存操作
	key := "name"
	value := "stormi-li"
	
	// 获取缓存数据
	res := cache.Get(key)
	if res == omicamo.NullString { // 如果缓存未命中
		cache.Set(key, value) // 设置缓存
	}

	// 再次获取缓存数据
	res = cache.Get(key)
	fmt.Println(res) // 打印缓存值
}
```