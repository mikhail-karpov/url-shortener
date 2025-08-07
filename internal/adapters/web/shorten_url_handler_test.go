package web

import (
	"encoding/json"
	"github.com/mikhail-karpov/url-shortener/internal/adapters/memory"
	"github.com/mikhail-karpov/url-shortener/internal/application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenURLHandler(t *testing.T) {

	var (
		cmdHandler = application.NewShortenURLCmdHandler(memory.NewRepository())
		handler    = ShortenURLHandler(cmdHandler)
	)

	t.Run("returns short url", func(t *testing.T) {

		payload := `{"url":"https://example.com"}`
		request := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(payload))
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response ShortURLResponse

		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
		assert.Equal(t, "https://example.com", response.OriginalURL)
		assert.NotNil(t, response.Alias)
		assert.NotNil(t, response.CreatedAt)
	})

	t.Run("returns bad request", func(t *testing.T) {

		payloads := []string{
			"",
			`{"url":""}`,
			`{"url":null}`,
		}

		for _, payload := range payloads {
			request := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(payload))
			request.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		}
	})
}
