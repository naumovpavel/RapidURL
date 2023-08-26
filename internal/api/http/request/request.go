package request

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func PrepareRequest[T any](r *http.Request) (*T, error) {
	req, err := DecodeRequest[T](r)
	if err != nil {
		return nil, err
	}
	if err = ValidateRequest[T](req); err != nil {
		return nil, err
	}

	return req, nil
}

func DecodeRequest[T any](r *http.Request) (*T, error) {
	var req T

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errors.New("request body is empty")
		} else {
			return nil, errors.New("invalid request body")
		}
	}

	return &req, nil
}

func ValidateRequest[T any](req *T) error {
	if err := validator.New().Struct(req); err != nil {
		return ValidationError(err.(validator.ValidationErrors))
	}

	return nil
}

func ValidationError(errs validator.ValidationErrors) error {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "alphanum":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must has only numbers and letters", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s isn't email address", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return errors.New(strings.Join(errMsgs, ", "))
}
