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

// ShortenURLHandler godoc
//
// @tags		URL
// @summary		Shorten URL
// @param		longUrl		body		web.ShortenURLRequest	true	"long URL"
// @success		200			{object}	web.ShortURLResponse
// @failure		400			{object}	web.ErrResponse
// @router		/shorten	[post]
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
