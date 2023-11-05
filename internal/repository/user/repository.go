package user

import (
	"context"
	"errors"

	"RapidURL/internal/entity"
)

var (
	ErrUserNotFound     = errors.New("user with this email not found")
	ErrUserAlreadyExist = errors.New("user eith this email already exist")
)

type Repository interface {
	SaveUser(ctx context.Context, user DTO) error
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
}
