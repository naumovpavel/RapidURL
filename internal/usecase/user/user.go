package user

import (
	"errors"
	"fmt"

	"RapidURL/internal/entity"
	"RapidURL/internal/lib/auth"
	"RapidURL/internal/lib/random"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	SaveUser(user entity.User) error
	FindUserByEmail(email string) (*entity.User, error)
}

type Usecase struct {
	userStorage Storage
}

func New(storage Storage) *Usecase {
	return &Usecase{
		userStorage: storage,
	}
}

func (u *Usecase) CreateUser(userDTO CreateUserDTO) error {
	const op = "usecase.user.CreateUser"

	salt := random.NewRandomString(10)
	pass, err := bcrypt.GenerateFromPassword([]byte(salt+userDTO.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return u.userStorage.SaveUser(entity.User{
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: string(pass),
		Salt:     salt,
	})
}

var (
	ErrUserNotFound  = errors.New("user with this email not found")
	ErrIncorrectPass = errors.New("incorrect password")
)

func (u *Usecase) LoginUser(userDTO LoginUserDTO) (string, error) {
	user, err := u.userStorage.FindUserByEmail(userDTO.Email)
	if err != nil {
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.Salt+userDTO.Password))
	if err != nil {
		return "", ErrIncorrectPass
	}

	return auth.CreateJWT(user.Id)
}
