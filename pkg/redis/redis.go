package redis_db

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type RedisDB interface {
	GetInfo(key string, res interface{}) error
	GetKey(key string) error
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
	DeletePPrefix(prefix string) error
}

type redisDB struct {
	idx int
}

var redisDBInstance RedisDB
var redisOnce sync.Once

func NewRedisDB() RedisDB {
	redisOnce.Do(func() {
		redisDBInstance = &redisDB{
			idx: 1,
		}
	})
	return redisDBInstance
}

func (db redisDB) GetInfo(key string, res interface{}) error {
	if res == nil {
		return errors.New("value invalid")
	}

	client := Client(db.idx)
	data, err := client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return err
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), &res)
}

func (db redisDB) GetKey(key string) error {
	client := Client(db.idx)
	if _, err := client.Get(context.Background(), key).Result(); err != nil {
		return err
	}
	return nil
}

func (db redisDB) Set(key string, value interface{}, expiration time.Duration) error {
	if value == nil {
		return errors.New("value invalid")
	}

	client := Client(db.idx)
	newValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.Set(context.Background(),
		key,
		newValue,
		expiration).Err()

}

func (db redisDB) Delete(key string) error {
	client := Client(db.idx)
	return client.Del(context.Background(), key).Err()
}

func (db redisDB) DeletePPrefix(prefix string) error {
	client := Client(db.idx)
	keys, err := client.Keys(context.Background(), prefix+"*").Result()

	if errors.Is(err, redis.Nil) {
		return nil
	}

	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := client.Del(context.Background(), key).Err(); err != nil {
			return err
		}
	}

	return nil
}
