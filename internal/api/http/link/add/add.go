package add

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"RapidURL/internal/api/http/request"
	"RapidURL/internal/api/http/response"
	"RapidURL/internal/entity"
	"RapidURL/internal/lib/logger/sl"
	storage "RapidURL/internal/repository/link"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Alias string `json:"alias,omitempty"`
	Url   string `json:"url" validate:"required,url"`
}

type Response struct {
	Alias string `json:"alias"`
}

type Saver interface {
	SaveLink(ctx context.Context, link entity.Link) (string, error)
}

func New(log *slog.Logger, sv Saver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.http.link.add"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		req, err := request.PrepareRequest[Request](r)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			response.Error(w, r, err, 400)
			return
		}

		reqUrl, err := url.Parse(req.Url)
		if err != nil {
			log.Error("fail to parse reqUrl", sl.Err(err))
			response.Error(w, r, err, 500)
			return
		}

		alias, err := sv.SaveLink(r.Context(), entity.Link{
			Alias: req.Alias,
			Url:   *reqUrl,
			User: entity.User{
				Id: r.Context().Value("userId").(int),
			},
		})

		if err != nil {
			handleSaveLinkError(log, w, r, err)
			return
		}

		render.JSON(w, r, Response{
			Alias: alias,
		})
	}
}

func handleSaveLinkError(log *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	log.Error("failed to save link", sl.Err(err))
	if errors.Is(err, storage.ErrAliasAlreadyExist) {
		response.Error(w, r, err, 403)
	} else {
		response.InternalError(w, r)
	}
}
