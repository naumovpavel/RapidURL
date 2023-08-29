package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"RapidURL/internal/api/http/link/add"
	"RapidURL/internal/api/http/link/redirect"
	"RapidURL/internal/api/http/middleware/auth"
	"RapidURL/internal/api/http/user/login"
	"RapidURL/internal/api/http/user/register"
	"RapidURL/internal/config"
	linkStorage "RapidURL/internal/storage/postgres/link"
	userStorage "RapidURL/internal/storage/postgres/user"
	"RapidURL/internal/usecase/link"
	"RapidURL/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := initLogger(cfg)
	r := initRouter(cfg, log)

	log.Info("Starting server...")

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Timeout,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go srv.ListenAndServe()
	<-sig

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	srv.Shutdown(ctx)
}

func initRouter(cfg *config.Config, log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	userStorage := userStorage.New(cfg.Postgres, log)
	userUsecase := user.New(userStorage)

	linkStorage := linkStorage.New(cfg.Postgres, log)
	linkUsecase := link.New(linkStorage)

	r.Post("/user/register", register.New(userUsecase, log))
	r.Post("/user/login", login.New(userUsecase, log))
	r.With(auth.New(log)).Post("/link", add.New(log, linkUsecase))
	r.Get("/{alias}", redirect.New(log, linkUsecase))

	return r
}

func initLogger(cfg *config.Config) *slog.Logger {
	var log *slog.Logger

	switch cfg.Env {
	case "LOCAL":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "PROD":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	}

	return log
}
