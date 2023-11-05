package link

import (
	"net/url"

	"RapidURL/internal/entity"
)

type DTO struct {
	Id     int
	Alias  string
	Url    string
	UserId int `db:"user_id"`
}

func FromEntity(entity entity.Link) DTO {
	return DTO{
		Id:     entity.Id,
		Alias:  entity.Alias,
		Url:    entity.Url.String(),
		UserId: entity.User.Id,
	}
}

func ToEntity(dto DTO) (entity.Link, error) {
	url, err := url.Parse(dto.Url)
	if err != nil {
		return entity.Link{}, err
	}
	return entity.Link{
		Id:    dto.Id,
		Alias: dto.Alias,
		Url:   *url,
		User: entity.User{
			Id: dto.UserId,
		},
	}, nil
}
