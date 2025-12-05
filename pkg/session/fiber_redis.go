package session

import (
	"app-ecommerce/config"
	redis_db "app-ecommerce/pkg/redis"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type fiberRedis struct {
}

func NewFiberRedis() fiber.Storage {
	return &fiberRedis{}
}

func (s *fiberRedis) Delete(key string) error {

	cfg := config.GetConfig()

	client := redis_db.Client(cfg.Web.RedisSessionIndex)
	ctx := context.Background()

	err := client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *fiberRedis) Reset() error {
	cfg := config.GetConfig()

	client := redis_db.Client(cfg.Web.RedisSessionIndex)
	ctx := context.Background()

	err := client.FlushAll(ctx).Err()
	if err != nil {
		return err
	}

	return nil

}

func (s *fiberRedis) Close() error {

	return nil
}

func (s *fiberRedis) Get(key string) ([]byte, error) {

	cfg := config.GetConfig()

	client := redis_db.Client(cfg.Web.RedisSessionIndex)
	ctx := context.Background()

	data, err := client.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

func (s *fiberRedis) Set(key string, data []byte, ttl time.Duration) error {

	cfg := config.GetConfig()

	client := redis_db.Client(cfg.Web.RedisSessionIndex)
	ctx := context.Background()

	err := client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
