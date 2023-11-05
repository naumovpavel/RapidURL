package user

import (
	"context"
	"errors"
	"fmt"

	"RapidURL/internal/entity"
	"RapidURL/internal/lib/auth"
	"RapidURL/internal/lib/random"
	repository "RapidURL/internal/repository/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound  = errors.New("user with this email not found")
	ErrIncorrectPass = errors.New("incorrect password")
)

type Usecase struct {
	repo repository.Repository
}

func New(Repo repository.Repository) *Usecase {
	return &Usecase{
		repo: Repo,
	}
}

func (u *Usecase) CreateUser(ctx context.Context, user entity.User) error {
	const op = "usecase.user.CreateUser"

	salt := random.NewRandomString(10)
	pass, err := bcrypt.GenerateFromPassword([]byte(salt+user.Password), bcrypt.DefaultCost)
	user.Salt = salt
	user.Password = string(pass)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return u.repo.SaveUser(ctx, repository.FromEntity(user))
}

func (u *Usecase) LoginUser(ctx context.Context, email string, pass string) (string, error) {
	user, err := u.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.Salt+pass))
	if err != nil {
		return "", ErrIncorrectPass
	}

	return auth.CreateJWT(user.Id)
}
