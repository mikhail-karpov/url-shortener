package web

import (
	"context"
	"github.com/mikhail-karpov/url-shortener/internal/application"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"net/http"
)

type CmdHandler interface {
	ShortenURL(ctx context.Context, cmd *application.ShortenURLCmd) (*domain.ShortURL, error)
}

func ShortenURLHandler(cmdHandler CmdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var request ShortenURLRequest
		err := readJSON(w, r, &request)
		if err != nil || request.URL == "" {
			writeBadRequest(w, "invalid url")
			return
		}

		cmd := &application.ShortenURLCmd{OriginalURL: request.URL}
		shortUrl, err := cmdHandler.ShortenURL(r.Context(), cmd)
		if err != nil {
			writeErr(w, err)
			return
		}

		response := &ShortURLResponse{
			OriginalURL: shortUrl.OriginalURL,
			Alias:       shortUrl.Alias,
			CreatedAt:   shortUrl.CreatedAt,
		}
		writeOK(w, response)
	}
}
