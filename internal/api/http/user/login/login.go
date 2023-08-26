package login

import (
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
	response.Response
	Jwt string `json:"jwt,omitempty"`
}

type loginer interface {
	LoginUser(userDTO user.LoginUserDTO) (string, error)
}

func New(login loginer, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.http.user.login.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		req, err := request.PrepareRequest[Request](r)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.Error(err))
			return
		}

		jwt, err := login.LoginUser(user.LoginUserDTO{
			Email:    req.Email,
			Password: req.Password,
		})

		if err != nil {
			log.Error("failed to login user", sl.Err(err))
			if errors.Is(err, user.ErrUserNotFound) {
				render.JSON(w, r, response.Error(err))
			} else {
				render.JSON(w, r, response.Error(errors.New("internal error")))
			}
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: jwt,
		})
		render.JSON(w, r, Response{
			Response: response.Ok(),
			Jwt:      jwt,
		})
	}
}
