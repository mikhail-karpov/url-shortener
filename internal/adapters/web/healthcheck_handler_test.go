package web

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {

	t.Run("returns status code 200 when healthy", func(t *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		HealthcheckHandler().ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
