package memory

import (
	"context"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRepository(t *testing.T) {

	var (
		repo        = NewRepository()
		expectedUrl = &domain.ShortURL{
			Alias:       "alias",
			OriginalURL: "https://google.com",
			CreatedAt:   time.Now(),
		}
	)

	t.Run("stores url", func(t *testing.T) {
		err := repo.Add(context.Background(), expectedUrl)
		require.NoError(t, err)

		actualUrl, err := repo.Get(context.Background(), "alias")
		require.NoError(t, err)
		assert.Equal(t, expectedUrl, actualUrl)
	})

	t.Run("returns error if url not found", func(t *testing.T) {

		_, err := repo.Get(context.Background(), "notfound")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
