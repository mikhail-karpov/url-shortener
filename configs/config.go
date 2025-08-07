package configs

import (
	"github.com/mikhail-karpov/url-shortener/internal/adapters/redis"
	"os"
	"strconv"
)

type Config struct {
	HTTP struct {
		Port int
	}
	Redis redis.Config
}

func InitConfig() Config {
	return Config{
		HTTP: struct{ Port int }{Port: getInt("SERVER_PORT", 8080)},
		Redis: redis.Config{
			Addr:     getString("REDIS_ADDR", "localhost:6379"),
			Password: getString("REDIS_PASSWORD", ""),
			DB:       getInt("REDIS_DB", 0),
		},
	}
}

func getString(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func getInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return valueAsInt
}
