package redis

import (
	"context"

	"github.com/go-redis/redis"

	"crud/config"
	"crud/storage"
)

type Cache struct {
	client *redis.Client
	user   *UserRepo
	order  *OrderRepo
}

func NewRedis(ctx context.Context, cfg config.Config) (storage.CacheI, error) {

	var client = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	err := client.Ping().Err()

	return &Cache{
		client: client,
		user:   NewUserRepo(client),
		order:  NewOrderRepo(client),
	}, err
}

func (c *Cache) CloseDB() {
	c.client.Close()
}

func (c *Cache) User() storage.UserCacheI {

	if c.user == nil {
		c.user = NewUserRepo(c.user.client)
	}

	return c.user
}

func (c *Cache) Order() storage.OrderCacheI {

	if c.order == nil {
		c.order = NewOrderRepo(c.order.client)
	}

	return c.order
}
