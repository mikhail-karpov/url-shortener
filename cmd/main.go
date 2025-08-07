package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	r "github.com/mikhail-karpov/url-shortener/internal/adapters/redis"
	"github.com/mikhail-karpov/url-shortener/internal/adapters/web"
	"github.com/mikhail-karpov/url-shortener/internal/application"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	repo, err := initRedisRepository()
	if err != nil {
		panic(err)
	}

	cmdHandler := application.NewShortenURLCmdHandler(repo)
	queryHandler := application.NewShortURLQueryHandler(repo)
	server := initHttpServer(cmdHandler, queryHandler)

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

func initRedisRepository() (*r.Repository, error) {

	cfg := r.Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	cache, err := r.NewCache(cfg)
	if err != nil {
		return nil, err
	}

	return r.NewRepository(cache, 0), nil
}

func initHttpServer(
	cmd *application.ShortenURLCmdHandler,
	query *application.ShortURLQueryHandler) *http.Server {

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/shorten", web.ShortenURLHandler(cmd))
		r.Get("/{alias}", web.RedirectURLHandler(query))
	})

	router.Route("/health", func(r chi.Router) {
		r.Get("/liveness", web.HealthcheckHandler())
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
