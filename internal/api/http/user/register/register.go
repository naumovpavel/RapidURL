package register

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"RapidURL/internal/api/http/request"
	"RapidURL/internal/api/http/response"
	"RapidURL/internal/entity"
	"RapidURL/internal/lib/logger/sl"
	repository "RapidURL/internal/repository/user"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Name     string `json:"name" validate:"required,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum"`
}

type Response struct {
	message string
}

type registerer interface {
	CreateUser(ctx context.Context, user entity.User) error
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
			response.Error(w, r, err, 400)
			return
		}

		err = reg.CreateUser(r.Context(), entity.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		})

		if err != nil {
			handleCreateUserFailure(w, r, log, err)
			return
		}

		render.JSON(w, r, &Response{message: "successfully registered"})
	}
}

func handleCreateUserFailure(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	log.Error("failed to create user", sl.Err(err))
	if errors.Is(err, repository.ErrUserAlreadyExist) {
		response.Error(w, r, err, 409)
	} else {
		response.InternalError(w, r)
	}
}
