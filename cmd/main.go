package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikhail-karpov/url-shortener/configs"
	_ "github.com/mikhail-karpov/url-shortener/docs"
	r "github.com/mikhail-karpov/url-shortener/internal/adapters/redis"
	"github.com/mikhail-karpov/url-shortener/internal/adapters/web"
	"github.com/mikhail-karpov/url-shortener/internal/application"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title			URL Shortener
// @version			1.0

// @license.name	Apache 2.0
// @license.url 	https://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath		/api/v1

// @accept			json
// @produce			json
func main() {

	cfg := configs.InitConfig()

	redisClient, err := r.NewClient(cfg.Redis)
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(redisClient)
	if err != nil {
		log.Fatal(err)
	}

	repo := r.NewRepository(redisClient, 0)
	cmdHandler := application.NewShortenURLCmdHandler(repo)
	queryHandler := application.NewShortURLQueryHandler(repo)
	server := initHttpServer(cfg, cmdHandler, queryHandler)

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	log.Printf("server listening on %s\n", server.Addr)

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-sigCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("unable to shutdown http server: %s\n", err)
	} else {
		log.Println("server gracefully stopped")
	}
}

func initHttpServer(
	cfg configs.Config,
	cmd *application.ShortenURLCmdHandler,
	query *application.ShortURLQueryHandler) *http.Server {

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/shorten", web.ShortenURLHandler(cmd))
		r.Get("/{id}", web.GetShortURLHandler(query))
	})

	router.Route("/", func(r chi.Router) {
		r.Get("/health/liveness", web.HealthcheckHandler())
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: router,
	}
}
