package web

import (
	"context"
	"net/http"

	"github.com/mikhail-karpov/url-shortener/internal/application"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
)

type CmdHandler interface {
	ShortenURL(ctx context.Context, cmd *application.ShortenURLCmd) (*domain.ShortURL, error)
}

func ShortenURLHandler(cmdHandler CmdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var request ShortenURLRequest
		err := readJSON(w, r, &request)
		if err != nil || request.LongURL == "" {
			writeBadRequest(w, "invalid url")
			return
		}

		cmd := &application.ShortenURLCmd{OriginalURL: request.LongURL}
		shortUrl, err := cmdHandler.ShortenURL(r.Context(), cmd)
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
