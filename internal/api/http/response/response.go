package response

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	internalError = errors.New("internal error")
)

type Response struct {
	Error string `json:"error,omitempty"`
}

func Error(w http.ResponseWriter, r *http.Request, err error, code int) {
	render.Status(r, code)
	render.JSON(w, r, &Response{Error: err.Error()})
}

func InternalError(w http.ResponseWriter, r *http.Request) {
	Error(w, r, internalError, 500)
}
