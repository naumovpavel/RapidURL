package main

import (
	"log/slog"
	"net/http"
	"os"

	"RapidURL/internal/api/http/link/add"
	"RapidURL/internal/api/http/link/redirect"
	"RapidURL/internal/api/http/middleware/auth"
	"RapidURL/internal/api/http/user/login"
	"RapidURL/internal/api/http/user/register"
	"RapidURL/internal/config"
	link2 "RapidURL/internal/storage/postgres/link"
	"RapidURL/internal/storage/postgres/user"
	"RapidURL/internal/usecase/link"
	user2 "RapidURL/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := initLogger(cfg)
	r := initRouter()
	s := user.New(cfg.Postgres, log)
	u := user2.New(s)

	ls := link2.New(cfg.Postgres, log)
	lu := link.New(ls)

	r.Post("/user/register", register.New(u, log))
	r.Post("/user/login", login.New(u, log))
	r.With(auth.New(log)).Post("/link", add.New(log, lu))
	r.Get("/{alias}", redirect.New(log, lu))

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
