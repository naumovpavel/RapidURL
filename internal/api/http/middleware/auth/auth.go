package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"RapidURL/internal/api/http/response"
	"RapidURL/internal/lib/auth"
	"RapidURL/internal/lib/logger/sl"
	"github.com/go-chi/render"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			jwt, err := r.Cookie("jwt")
			if err != nil {
				unAuth(log, err, w, r)
				return
			}

			userId, err := auth.DecodeJWT(jwt.Value)
			if err != nil {
				unAuth(log, err, w, r)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func unAuth(log *slog.Logger, err error, w http.ResponseWriter, r *http.Request) {
	log.Error("user unauthorized", sl.Err(err))
	render.JSON(w, r, response.Error(errors.New("unauthorized")))
}
