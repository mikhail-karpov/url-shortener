package redis

import (
	"context"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestRedisRepository(t *testing.T) {

	ctx := context.Background()
	request := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7.4.1-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, request)
	defer testcontainers.CleanupContainer(t, container)
	require.NoError(t, err)

	redisAddr, err := container.Endpoint(ctx, "")
	require.NoError(t, err)

	client, err := NewClient(Config{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	defer client.Close()
	require.NoError(t, err)

	repository := NewRepository(client, 0)

	t.Run("stores url", func(t *testing.T) {

		expected := &domain.ShortURL{
			ID:        "google",
			LongURL:   "https://google.com",
			CreatedAt: time.Now(),
		}

		err := repository.Add(ctx, expected)
		require.NoError(t, err)

		actual, err := repository.Get(ctx, "google")
		require.NoError(t, err)
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.LongURL, actual.LongURL)
	})

	t.Run("returns error if url not found", func(t *testing.T) {

		_, err := repository.Get(ctx, "not-found")
		assert.Equal(t, err, domain.ErrNotFound)
	})
}
