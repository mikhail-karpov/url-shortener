package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikhail-karpov/url-shortener/internal/adapters/memory"
	"github.com/mikhail-karpov/url-shortener/internal/application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURLHandler(t *testing.T) {

	var (
		cmdHandler = application.NewShortenURLCmdHandler(memory.NewRepository())
		handler    = ShortenURLHandler(cmdHandler)
	)

	t.Run("returns short url", func(t *testing.T) {

		payload := `{"long_url":"https://example.com"}`
		request := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(payload))
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response ShortURLResponse
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
		assert.Equal(t, "https://example.com", response.LongURL)
		assert.NotNil(t, response.ID)
		assert.NotNil(t, response.CreatedAt)
	})

	t.Run("returns bad request", func(t *testing.T) {

		payloads := []string{
			"",
			`{"long_url":""}`,
			`{"long_url":null}`,
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
