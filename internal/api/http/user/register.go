package user

import (
	"context"
	"log/slog"

	"RapidURL/internal/api/http/user/login"
	"RapidURL/internal/api/http/user/register"
	"RapidURL/internal/entity"
	"github.com/go-chi/chi/v5"
)

type usecase interface {
	CreateUser(ctx context.Context, user entity.User) error
	LoginUser(ctx context.Context, email string, pass string) (string, error)
}

func Register(r *chi.Mux, userUsecase usecase, log *slog.Logger) {
	r.Post("/user/register", register.New(userUsecase, log))
	r.Post("/user/login", login.New(userUsecase, log))
}
