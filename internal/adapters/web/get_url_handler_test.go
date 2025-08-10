package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mikhail-karpov/url-shortener/internal/adapters/memory"
	"github.com/mikhail-karpov/url-shortener/internal/application"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectURLHandler(t *testing.T) {

	var (
		repo         = memory.NewRepository()
		queryHandler = application.NewShortURLQueryHandler(repo)
		handler      = GetShortURLHandler(queryHandler)
	)

	t.Run("returns 200 with short url", func(t *testing.T) {

		err := repo.Add(context.Background(), &domain.ShortURL{
			LongURL:   "https://example.com",
			ID:        "example",
			CreatedAt: time.Now(),
		})
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.SetPathValue("id", "example")

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		var response ShortURLResponse
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
		assert.Equal(t, "https://example.com", response.LongURL)
		assert.Equal(t, "example", response.ID)
		assert.NotNil(t, response.CreatedAt)
	})

	t.Run("returns 404 if url not found", func(t *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.SetPathValue("id", "google")

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
}
