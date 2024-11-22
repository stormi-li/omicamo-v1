package omicamo

import "github.com/go-redis/redis/v8"

type Client struct {
	redisClient *redis.Client
}

func NewClient(opts *redis.Options) *Client {
	return &Client{redisClient: redis.NewClient(opts)}
}

func (c *Client) NewCache() *Cache {
	return newCache(c.redisClient)
}
