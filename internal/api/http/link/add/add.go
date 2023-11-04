package add

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"RapidURL/internal/api/http/request"
	"RapidURL/internal/api/http/response"
	"RapidURL/internal/lib/logger/sl"
	link2 "RapidURL/internal/storage/postgres/link"
	"RapidURL/internal/usecase/link"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Alias string `json:"alias,omitempty"`
	Url   string `json:"url" validate:"required,url"`
}

type Response struct {
	response.Response
	Alias string `json:"alias"`
}

type Saver interface {
	SaveLink(link link.SaveLinkDTO) (string, error)
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
			handleInvalidRequestError(log, w, r, err)
			return
		}

		reqUrl, err := url.Parse(req.Url)
		if err != nil {
			handleParseUrlError(w, r, log, err)
			return
		}

		alias, err := sv.SaveLink(link.SaveLinkDTO{
			Alias:  req.Alias,
			Url:    *reqUrl,
			UserId: r.Context().Value("userId").(int),
		})

		if err != nil {
			handleSaveLinkError(log, w, r, err)
			return
		}

		render.JSON(w, r, Response{
			Response: response.Ok(),
			Alias:    alias,
		})
	}
}

func handleParseUrlError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	log.Error("fail to parse reqUrl", sl.Err(err))
	render.Status(r, 500)
	render.JSON(w, r, response.Error(errors.New("internal error")))
}

func handleInvalidRequestError(log *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	log.Error("invalid request", sl.Err(err))
	render.Status(r, 400)
	render.JSON(w, r, response.Error(err))
}

func handleSaveLinkError(log *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	log.Error("failed to save link", sl.Err(err))
	if errors.Is(err, link2.ErrAliasAlreadyExist) {
		render.Status(r, 403)
		render.JSON(w, r, response.Error(err))
	} else {
		render.Status(r, 500)
		render.JSON(w, r, response.Error(errors.New("internal error")))
	}
}
