package redis

import (
	"context"
	"fmt"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"time"
)

type redisShortUrl struct {
	OriginalUrl string    `redis:"url"`
	CreatedAt   time.Time `redis:"created_at"`
}

type Repository struct {
	cache *Cache
	ttl   time.Duration
}

func NewRepository(cache *Cache, ttl time.Duration) *Repository {
	return &Repository{
		cache: cache,
		ttl:   ttl,
	}
}

func (r *Repository) Add(ctx context.Context, shortUrl *domain.ShortURL) error {

	key := key(shortUrl.Alias)
	value := &redisShortUrl{
		OriginalUrl: shortUrl.OriginalURL,
		CreatedAt:   shortUrl.CreatedAt,
	}

	err := r.cache.Put(ctx, key, value, r.ttl)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, alias string) (*domain.ShortURL, error) {

	key := key(alias)
	var redisShortUrl redisShortUrl
	err := r.cache.Get(ctx, key, &redisShortUrl)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	shortUrl := &domain.ShortURL{
		Alias:       alias,
		OriginalURL: redisShortUrl.OriginalUrl,
		CreatedAt:   redisShortUrl.CreatedAt,
	}
	return shortUrl, nil
}

func key(alias string) string {
	return fmt.Sprintf("url-shortener:urls:%s", alias)
}
