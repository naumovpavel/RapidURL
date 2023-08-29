package link

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"RapidURL/internal/config"
	"RapidURL/internal/entity"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func New(cfg config.Postgres, log *slog.Logger) *Storage {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Database,
	)

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Error("can't open database connection", err)
	}

	err = db.Ping()
	if err != nil {
		log.Error("database doesn't response", err)
	}

	return &Storage{
		db:  db,
		log: log,
	}
}

var ErrAliasAlreadyExist = errors.New("this alias already exist")

func (s *Storage) SaveLink(link *entity.Link) error {
	const op = "storage.postgres.SaveLink"

	stmt, err := s.db.Prepare("insert into links(alias, url, user_id) VALUES ($1,$2,$3)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(link.Alias, link.Url.String(), link.UserId)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAliasAlreadyExist
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

var ErrLinkNotFound = errors.New("link with this alias not found")

func (s *Storage) FindLinkByAlias(alias string) (*entity.Link, error) {
	const op = "storage.postgres.FindUserByEmail"

	stmt, err := s.db.Prepare("select alias, url, user_id from links where alias = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var link entity.Link
	row := stmt.QueryRow(alias)
	err = row.Scan(&link.Alias, &link.Url, &link.UserId)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "02000" {
			return nil, ErrLinkNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &link, nil
}
