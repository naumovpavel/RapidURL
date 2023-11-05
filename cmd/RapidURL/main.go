package main

import (
	"context"
	"fmt"
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
	linkRepository "RapidURL/internal/repository/link"
	memcachedLinkRepository "RapidURL/internal/repository/link/memcached"
	postgresLinkRepository "RapidURL/internal/repository/link/postgres"
	userRepository "RapidURL/internal/repository/user/postgres"
	"RapidURL/internal/usecase/link"
	"RapidURL/internal/usecase/user"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()
	log := initLogger(cfg)
	pool := initPool(cfg.Postgres, log)
	defer pool.Close()
	memcache := memcache.New("localhost:11211")
	defer memcache.Close()
	r := initRouter(cfg, log, pool, memcache)

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

func initRouter(cfg *config.Config, log *slog.Logger, pool *pgxpool.Pool, memcache *memcache.Client) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	userStorage := userRepository.New(pool, log)
	userUsecase := user.New(userStorage)

	linkStorage := postgresLinkRepository.New(pool, log)
	linkCache := memcachedLinkRepository.New(memcache)
	cachedLink := linkRepository.NewCachedRepository(linkStorage, linkCache, log)
	linkUsecase := link.New(cachedLink)

	r.Post("/user/register", register.New(userUsecase, log))
	r.Post("/user/login", login.New(userUsecase, log))
	r.With(auth.New(log)).Post("/link/add", add.New(log, linkUsecase))
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

func initPool(cfg config.Postgres, log *slog.Logger) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	connUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Database,
	)

	pool, err := pgxpool.New(ctx, connUrl)
	if err != nil {
		log.Error("can't open database connection", err)
		panic(err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Error("database doesn't response", err)
		panic(err)
	}

	return pool
}
