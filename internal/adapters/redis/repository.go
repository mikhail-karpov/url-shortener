package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"github.com/redis/go-redis/v9"
)

type redisShortUrl struct {
	LongURL   string    `redis:"url"`
	CreatedAt time.Time `redis:"created_at"`
}

type Repository struct {
	cache *Cache
	ttl   time.Duration
}

func NewRepository(client *redis.Client, ttl time.Duration) *Repository {
	return &Repository{
		cache: &Cache{client: client},
		ttl:   ttl,
	}
}

func (r *Repository) Add(ctx context.Context, shortUrl *domain.ShortURL) error {

	key := key(shortUrl.ID)
	value := &redisShortUrl{
		LongURL:   shortUrl.LongURL,
		CreatedAt: shortUrl.CreatedAt,
	}

	err := r.cache.Put(ctx, key, value, r.ttl)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*domain.ShortURL, error) {

	key := key(id)
	var redisShortUrl redisShortUrl
	err := r.cache.Get(ctx, key, &redisShortUrl)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	shortUrl := &domain.ShortURL{
		ID:        id,
		LongURL:   redisShortUrl.LongURL,
		CreatedAt: redisShortUrl.CreatedAt,
	}
	return shortUrl, nil
}

func key(alias string) string {
	return fmt.Sprintf("url-shortener:urls:%s", alias)
}
