package omicamo

import "github.com/go-redis/redis/v8"

type Client struct {
	RedisClient *redis.Client
}

func NewClient(opts *redis.Options) *Client {
	return &Client{RedisClient: redis.NewClient(opts)}
}

func (c *Client) NewCache() *Cache {
	return newCache(c.RedisClient)
}
