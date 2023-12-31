package login

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"RapidURL/internal/api/http/request"
	"RapidURL/internal/api/http/response"
	"RapidURL/internal/lib/logger/sl"
	"RapidURL/internal/usecase/user"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum"`
}

type Response struct {
	Jwt string `json:"jwt,omitempty"`
}

type loginer interface {
	LoginUser(ctx context.Context, email string, pass string) (string, error)
}

func New(login loginer, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.http.user.login"
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

		jwt, err := login.LoginUser(r.Context(), req.Email, req.Password)

		if err != nil {
			handleLoginFailure(w, r, log, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: jwt,
			Path:  "/",
		})
		render.JSON(w, r, Response{
			Jwt: jwt,
		})
	}
}

func handleLoginFailure(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	log.Error("failed to login user", sl.Err(err))
	if errors.Is(err, user.ErrUserNotFound) {
		response.Error(w, r, err, 404)
	} else if errors.Is(err, user.ErrIncorrectPass) {
		response.Error(w, r, err, 401)
	} else {
		response.InternalError(w, r)
	}
}
