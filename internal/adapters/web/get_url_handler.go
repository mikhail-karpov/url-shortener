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

// GetShortURLHandler godoc
//
// @tags 		URL
// @summary 	Get URL
// @param		id		path		string true	"URL id"
// @success 	200		{object}	web.ShortURLResponse
// @failure 	404		{object}	web.ErrResponse
// @router		/{id}	[get]
func GetShortURLHandler(provider ShortURLProvider) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		if id == "" {
			writeBadRequest(w, "invalid alias")
			return
		}

		query := application.ShortURLQuery{ID: id}
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
	}
}
