package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

type Cache struct {
	client *redis.Client
}

func NewCache(cfg Config) (*Cache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	if pong != "PONG" {
		return nil, fmt.Errorf("failed to ping redis server")
	}

	return &Cache{client: client}, nil
}

func (cache *Cache) Put(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal data: %s", err)
		return err
	}

	_, err = cache.client.Set(ctx, key, value, ttl).Result()
	if err != nil {
		log.Printf("failed to store data: %s", err)
		return err
	}
	return nil
}

func (cache *Cache) Get(ctx context.Context, key string, data interface{}) error {
	value, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("failed to get data: %s", err)
		return err
	}

	err = json.Unmarshal([]byte(value), data)
	if err != nil {
		log.Printf("failed to unmarshal json: %s", err)
		return err
	}
	return nil
}
