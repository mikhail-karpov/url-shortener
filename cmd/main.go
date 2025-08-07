package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikhail-karpov/url-shortener/internal/adapters/memory"
	"github.com/mikhail-karpov/url-shortener/internal/adapters/web"
	"github.com/mikhail-karpov/url-shortener/internal/application"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	repo := memory.NewRepository()
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

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("unable to shutdown http server: %s\n", err)
	} else {
		log.Println("server gracefully stopped")
	}
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

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
