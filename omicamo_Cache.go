package omicamo

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omisync-v1"
)

// 分布式锁前缀
var DLockPrefix = "stormi:dlock:"
var NullString = "_null"

// Cache 结构体定义
type Cache struct {
	RedisClient         *redis.Client        // Redis 客户端
	ctx                 context.Context      // 上下文，用于 Redis 操作
	omisyncClient       *omisync.Client      // Omisync 分布式锁客户端
	cacheGetCallback    func(string) string  // 缓存读取回调
	cacheSetCallback    func(string, string) // 缓存写入回调
	databaseGetCallback func(string) string  // 数据库读取回调
	databaseSetCallback func(string, string) // 数据库写入回调

}

// NewCache 创建一个新的 Cache 实例
func newCache(redisClient *redis.Client) *Cache {
	return &Cache{
		RedisClient:   redisClient,
		ctx:           context.Background(),
		omisyncClient: omisync.NewClient(redisClient.Options()),
	}
}

// 设置回调函数
func (c *Cache) AddCallback(cacheGet, databaseGet func(key string) string, cacheSet, databaseSet func(key string, value string)) {
	c.cacheGetCallback = cacheGet
	c.cacheSetCallback = cacheSet
	c.databaseGetCallback = databaseGet
	c.databaseSetCallback = databaseSet
}

// Get 从缓存中读取值，如果缓存中不存在，则从数据库获取
func (c *Cache) Get(key string) string {
	// 从缓存中获取值
	res := c.cacheGetCallback(key)
	if res != "" {
		return res
	}

	// 加分布式锁
	dLock := c.omisyncClient.NewLock(DLockPrefix + key)
	dLock.Lock()
	defer dLock.Unlock()

	res = c.cacheGetCallback(key)
	if res != "" {
		return res
	}

	// 从数据库中获取值
	res = c.databaseGetCallback(key)
	if res == "" {
		res = NullString
	}
	c.cacheSetCallback(key, res)
	return res
}

// Set 更新缓存和数据库的值
func (c *Cache) Set(key, value string) {
	// 加锁
	dLock := c.omisyncClient.NewLock(DLockPrefix + key)
	dLock.Lock()
	defer dLock.Unlock()

	c.RedisClient.Del(c.ctx, key)
	// 更新数据
	c.databaseSetCallback(key, value)
	c.cacheSetCallback(key, value)
}
