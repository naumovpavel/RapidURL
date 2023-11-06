package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"RapidURL/internal/app"
	"RapidURL/internal/config"
	"RapidURL/internal/repository/postgres"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	cfg := config.MustLoad()
	log := initLogger(cfg)
	pool := postgres.InitPool(cfg.Postgres, log)
	defer pool.Close()
	memcache := memcache.New(cfg.Memcached.Hosts...)
	defer memcache.Close()

	app := app.New(pool, memcache, log, cfg)
	app.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	app.Stop(ctx)
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
