package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"RapidURL/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPool(cfg config.Postgres, log *slog.Logger) *pgxpool.Pool {
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
