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
	r "github.com/mikhail-karpov/url-shortener/internal/adapters/redis"
	"github.com/mikhail-karpov/url-shortener/internal/adapters/web"
	"github.com/mikhail-karpov/url-shortener/internal/application"
	"github.com/redis/go-redis/v9"
)

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
		r.Get("/{alias}", web.GetShortURLHandler(query))
	})

	router.Route("/health", func(r chi.Router) {
		r.Get("/liveness", web.HealthcheckHandler())
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: router,
	}
}
