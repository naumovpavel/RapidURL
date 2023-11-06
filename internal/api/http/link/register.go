package link

import (
	"context"
	"log/slog"
	"net/url"

	"RapidURL/internal/api/http/link/add"
	"RapidURL/internal/api/http/link/redirect"
	"RapidURL/internal/api/http/middleware/auth"
	"RapidURL/internal/entity"
	"github.com/go-chi/chi/v5"
)

type usecase interface {
	SaveLink(ctx context.Context, link entity.Link) (string, error)
	GetLink(ctx context.Context, alias string) (url.URL, error)
}

func Register(r *chi.Mux, linkUsecase usecase, log *slog.Logger) {
	r.With(auth.New(log)).Post("/link/add", add.New(log, linkUsecase))
	r.Get("/{alias}", redirect.New(log, linkUsecase))
}
