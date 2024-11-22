package omicamo

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omisync-v1"
)

// 分布式锁前缀
var DLockPrefix = "stormi:dlock:"
var NullString = "_null"

// Cache 结构体定义
type Cache struct {
	redisClient      *redis.Client                       // Redis 客户端
	ctx              context.Context                     // 上下文，用于 Redis 操作
	omisyncClient    *omisync.Client                     // Omisync 分布式锁客户端
	DelayDuration    time.Duration                       // 延迟时间，用于延时双删
	cacheGetCallback func(string, *redis.Client) string  // 缓存读取回调
	cacheSetCallback func(string, string, *redis.Client) // 缓存写入回调
	dbGetCallback    func(string) string                 // 数据库读取回调
	dbSetCallback    func(string, string)                // 数据库写入回调

}

// NewCache 创建一个新的 Cache 实例
func newCache(redisClient *redis.Client) *Cache {
	return &Cache{
		redisClient:   redisClient,
		ctx:           context.Background(),
		omisyncClient: omisync.NewClient(redisClient.Options()),
		DelayDuration: 2 * time.Second, // 默认延迟 2 秒
	}
}

// SetCacheCallback 设置缓存的回调函数
func (c *Cache) SetCacheCallback(get func(key string, redisClient *redis.Client) string,
	set func(key string, value string, redisClient *redis.Client)) {
	c.cacheGetCallback = get
	c.cacheSetCallback = set
}

// SetDatabaseCallback 设置数据库的回调函数
func (c *Cache) SetDatabaseCallback(get func(key string) string,
	set func(key string, value string)) {
	c.dbGetCallback = get
	c.dbSetCallback = set
}

// Get 从缓存中读取值，如果缓存中不存在，则从数据库获取
func (c *Cache) Get(key string) string {
	if c.cacheGetCallback == nil || c.dbGetCallback == nil {
		panic("Get callback is not set")
	}
	// 从缓存中获取值
	res := c.cacheGetCallback(key, c.redisClient)
	if res != "" {
		return res
	}

	// 加分布式锁
	dLock := c.omisyncClient.NewLock(DLockPrefix + key)
	dLock.Lock()
	defer dLock.Unlock()

	res = c.cacheGetCallback(key, c.redisClient)
	if res != "" {
		return res
	}

	// 从数据库中获取值
	res = c.dbGetCallback(key)
	if res == "" {
		res = NullString
	}
	c.cacheSetCallback(key, res, c.redisClient)
	return res
}

// Set 更新缓存和数据库的值
func (c *Cache) Set(key, value string) {
	if c.dbSetCallback == nil || c.cacheSetCallback == nil {
		panic("Set callback is not set")
	}

	// 加锁
	dLock := c.omisyncClient.NewLock(DLockPrefix + key)
	dLock.Lock()
	defer dLock.Unlock()

	c.redisClient.Del(c.ctx, key)
	// 更新数据
	c.dbSetCallback(key, value)
	c.cacheSetCallback(key, value, c.redisClient)
}
