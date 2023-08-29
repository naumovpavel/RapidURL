package register

import (
	"errors"
	"log/slog"
	"net/http"

	"RapidURL/internal/api/http/request"
	"RapidURL/internal/api/http/response"
	"RapidURL/internal/lib/logger/sl"
	user2 "RapidURL/internal/storage/postgres/user"
	"RapidURL/internal/usecase/user"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Name     string `json:"name" validate:"required,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum"`
}

type Response struct {
	response.Response
}

type registerer interface {
	CreateUser(userDTO user.CreateUserDTO) error
}

func New(reg registerer, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.http.user.register.New"
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

		err = reg.CreateUser(user.CreateUserDTO{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		})

		if err != nil {
			log.Error("failed to create user", sl.Err(err))
			if errors.Is(err, user2.ErrUserAlreadyExist) {
				render.JSON(w, r, response.Error(err))
			} else {
				render.JSON(w, r, response.Error(errors.New("internal error")))
			}
			return
		}

		render.JSON(w, r, response.Ok())
	}
}
