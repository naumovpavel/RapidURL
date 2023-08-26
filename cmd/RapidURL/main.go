package main

import (
	"log/slog"
	"net/http"
	"os"

	"RapidURL/internal/api/http/links/add"
	"RapidURL/internal/api/http/middleware/auth"
	"RapidURL/internal/api/http/user/login"
	"RapidURL/internal/api/http/user/register"
	"RapidURL/internal/config"
	"RapidURL/internal/storage/postgres"
	user2 "RapidURL/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := initLogger(cfg)
	r := initRouter()
	s := postgres.NewUserStorage(cfg.Postgres, log)
	u := user2.New(s)

	r.Post("/user/register", register.New(u, log))
	r.Post("/user/login", login.New(u, log))
	r.With(auth.New(log)).Get("/link", add.New(log))

	log.Info("Starting server...")

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Timeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}
}

func initRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

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
