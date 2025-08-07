package application

import (
	"context"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
)

type ShortURLQuery struct {
	Alias string
}

type ShortURLProvider interface {
	Get(ctx context.Context, query string) (*domain.ShortURL, error)
}

type ShortURLQueryHandler struct {
	provider ShortURLProvider
}

func NewShortURLQueryHandler(provider ShortURLProvider) *ShortURLQueryHandler {
	return &ShortURLQueryHandler{provider: provider}
}

func (h *ShortURLQueryHandler) Get(
	ctx context.Context, query ShortURLQuery) (*domain.ShortURL, error) {

	shortUrl, err := h.provider.Get(ctx, query.Alias)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return shortUrl, nil
}
