package memory

import (
	"context"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"sync"
)

type Repository struct {
	m    sync.RWMutex
	urls map[string]domain.ShortURL
}

func NewRepository() *Repository {
	return &Repository{
		m:    sync.RWMutex{},
		urls: make(map[string]domain.ShortURL),
	}
}

func (r *Repository) Add(_ context.Context, shortUrl *domain.ShortURL) error {
	r.m.Lock()
	defer r.m.Unlock()

	r.urls[shortUrl.Alias] = *shortUrl
	return nil
}

func (r *Repository) Get(_ context.Context, alias string) (*domain.ShortURL, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	shortUrl, ok := r.urls[alias]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &shortUrl, nil
}
