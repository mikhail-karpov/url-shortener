package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"time"
)

const (
	alphabet    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	alphabetLen = uint32(len(alphabet))
)

type ShortenURLCmd struct {
	OriginalURL string
}

type Storage interface {
	Add(ctx context.Context, shortUrl *domain.ShortURL) error
}

type ShortenURLCmdHandler struct {
	storage Storage
}

func NewShortenURLCmdHandler(storage Storage) *ShortenURLCmdHandler {
	if storage == nil {
		panic("storage is nil")
	}
	return &ShortenURLCmdHandler{storage: storage}
}

func (h *ShortenURLCmdHandler) ShortenURL(
	ctx context.Context, cmd *ShortenURLCmd) (*domain.ShortURL, error) {

	id, err := uuid.NewUUID()
	if err != nil {
		return &domain.ShortURL{}, err
	}

	shortUrl := &domain.ShortURL{
		LongURL:   cmd.OriginalURL,
		ID:        shorten(id.ID()),
		CreatedAt: time.Now(),
	}

	err = h.storage.Add(ctx, shortUrl)
	if err != nil {
		return &domain.ShortURL{}, err
	}
	return shortUrl, nil
}
