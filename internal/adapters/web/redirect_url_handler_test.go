package web

import (
	"context"
	"github.com/mikhail-karpov/url-shortener/internal/adapters/memory"
	"github.com/mikhail-karpov/url-shortener/internal/application"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRedirectURLHandler(t *testing.T) {

	var (
		repo         = memory.NewRepository()
		queryHandler = application.NewShortURLQueryHandler(repo)
		handler      = RedirectURLHandler(queryHandler)
	)

	t.Run("redirects to original url", func(t *testing.T) {

		err := repo.Add(context.Background(), &domain.ShortURL{
			OriginalURL: "https://example.com",
			Alias:       "example",
			CreatedAt:   time.Now(),
		})
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.SetPathValue("alias", "example")

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusFound, recorder.Code)
		assert.Equal(t, "https://example.com", recorder.Header().Get("Location"))
	})

	t.Run("returns 404 if url not found", func(t *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.SetPathValue("alias", "google")

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
}
