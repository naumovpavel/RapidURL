package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"RapidURL/internal/entity"
	repository "RapidURL/internal/repository/user"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
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

func (s *Repository) SaveUser(ctx context.Context, user repository.DTO) error {
	const op = "repository.postgres.SaveUser"

	_, err := s.pool.Exec(ctx, "insert into users(name, password, email, salt) VALUES ($1,$2,$3,$4)",
		user.Name,
		user.Password,
		user.Email,
		user.Salt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repository.ErrUserAlreadyExist
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Repository) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
	const op = "repository.postgres.FindUserByEmail"

	row := s.pool.QueryRow(ctx, "select users.id, users.email, users.name, users.salt, users.password from users where users.email = $1", email)
	var user repository.DTO
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.Salt, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, repository.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return repository.ToEntity(user), nil
}
