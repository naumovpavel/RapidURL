package link

import (
	"context"
	"errors"

	"RapidURL/internal/entity"
)

var (
	ErrAliasAlreadyExist = errors.New("this alias already exist")
	ErrLinkNotFound      = errors.New("link with this alias not found")
)

type Repository interface {
	SaveLink(ctx context.Context, link DTO) error
	FindLinkByAlias(ctx context.Context, alias string) (entity.Link, error)
}
