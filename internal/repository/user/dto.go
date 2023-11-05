package user

import (
	"RapidURL/internal/entity"
)

type DTO struct {
	Id       int
	Name     string
	Password string
	Salt     string
	Email    string
}

func FromEntity(user entity.User) DTO {
	return DTO{
		Id:       user.Id,
		Email:    user.Email,
		Name:     user.Name,
		Salt:     user.Salt,
		Password: user.Password,
	}
}

func ToEntity(user DTO) entity.User {
	return entity.User{
		Id:       user.Id,
		Email:    user.Email,
		Name:     user.Name,
		Salt:     user.Salt,
		Password: user.Password,
	}
}
