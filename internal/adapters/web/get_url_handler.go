package web

import (
	"context"
	"net/http"

	"github.com/mikhail-karpov/url-shortener/internal/application"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
)

type ShortURLProvider interface {
	Get(ctx context.Context, query application.ShortURLQuery) (*domain.ShortURL, error)
}

func GetShortURLHandler(provider ShortURLProvider) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		alias := r.PathValue("alias")
		if alias == "" {
			writeBadRequest(w, "invalid alias")
			return
		}

		query := application.ShortURLQuery{ID: alias}
		shortUrl, err := provider.Get(r.Context(), query)
		if err != nil {
			writeErr(w, err)
			return
		}

		response := &ShortURLResponse{
			ID:        shortUrl.ID,
			LongURL:   shortUrl.LongURL,
			CreatedAt: shortUrl.CreatedAt,
		}
		writeOK(w, response)
	})
}
