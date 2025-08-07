package web

import (
	"context"
	"github.com/mikhail-karpov/url-shortener/internal/application"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"net/http"
)

type ShortURLProvider interface {
	Get(ctx context.Context, query application.ShortURLQuery) (*domain.ShortURL, error)
}

func RedirectURLHandler(provider ShortURLProvider) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		alias := r.PathValue("alias")
		if alias == "" {
			writeBadRequest(w, "invalid alias")
			return
		}

		query := application.ShortURLQuery{Alias: alias}
		shortUrl, err := provider.Get(r.Context(), query)
		if err != nil {
			writeErr(w, err)
			return
		}

		http.Redirect(w, r, shortUrl.OriginalURL, http.StatusFound)
	})
}
