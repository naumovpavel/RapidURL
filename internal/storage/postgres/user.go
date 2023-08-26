package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"RapidURL/internal/config"
	"RapidURL/internal/entity"
	"RapidURL/internal/lib/logger/sl"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type UserStorage struct {
	db  *sql.DB
	log *slog.Logger
}

var (
	userExist = errors.New("user with this email already exist")
)

func NewUserStorage(cfg config.Postgres, log *slog.Logger) *UserStorage {
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

	return &UserStorage{
		db:  db,
		log: log,
	}
}

var ErrUserAlreadyExist = errors.New("user eith this email already exist")

func (s *UserStorage) SaveUser(user entity.User) error {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("insert into users(name, password, email, salt) VALUES ($1,$2,$3,$4)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(user.Name, user.Password, user.Email, user.Salt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserAlreadyExist
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

var ErrUserNotFound = errors.New("user with this email not found")

func (s *UserStorage) FindUserByEmail(email string) (*entity.User, error) {
	const op = "storage.postgres.FindUserByEmail"

	stmt, err := s.db.Prepare("select users.id, users.email, users.name, users.salt, users.password from users where users.email = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user entity.User
	row := stmt.QueryRow(email)
	err = row.Scan(&user.Id, &user.Email, &user.Name, &user.Salt, &user.Password)

	if err != nil {
		var pgErr *pq.Error
		s.log.Error("fuck", sl.Err(err))
		if errors.As(err, &pgErr) && pgErr.Code == "02000" {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
