package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"RapidURL/internal/api/http/response"
	"RapidURL/internal/lib/logger/sl"
	storage "RapidURL/internal/storage/postgres/link"
	"RapidURL/internal/usecase/link"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Getter interface {
	GetLink(dto link.GetLinkDTO) (url.URL, error)
}

func New(log *slog.Logger, gt Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.http.link.redirect"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		redirectUrl, err := gt.GetLink(link.GetLinkDTO{Alias: alias})

		if err != nil {
			handleLinkFindError(w, r, log, err)
			return
		}

		http.Redirect(w, r, redirectUrl.String(), http.StatusTemporaryRedirect)
	}
}

func handleLinkFindError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	log.Error("failed to find link", sl.Err(err))
	if errors.Is(err, storage.ErrLinkNotFound) {
		response.Error(w, r, err, 404)
	} else {
		response.InternalError(w, r)
	}
}
