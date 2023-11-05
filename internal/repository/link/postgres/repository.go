package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"RapidURL/internal/entity"
	repository "RapidURL/internal/repository/link"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

var _ repository.Repository = &Repository{}

type Repository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func New(pool *pgxpool.Pool, log *slog.Logger) *Repository {
	return &Repository{
		pool: pool,
		log:  log,
	}
}

func (s *Repository) SaveLink(ctx context.Context, link repository.DTO) error {
	const op = "repository.postgres.SaveLink"

	_, err := s.pool.Exec(ctx, "insert into links(alias, url, user_id) VALUES ($1,$2,$3)",
		link.Alias,
		link.Url,
		link.UserId,
	)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repository.ErrAliasAlreadyExist
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Repository) FindLinkByAlias(ctx context.Context, alias string) (entity.Link, error) {
	const op = "repository.postgres.FindLinkByAlias"

	row := s.pool.QueryRow(ctx, "select alias, url, user_id from links where alias = $1", alias)
	var link repository.DTO

	err := row.Scan(&link.Alias, &link.Url, &link.UserId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Link{}, repository.ErrLinkNotFound
		}
		return entity.Link{}, fmt.Errorf("%s: %w", op, err)
	}

	return repository.ToEntity(link)
}
