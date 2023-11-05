package link

import (
	"context"
	"net/url"

	"RapidURL/internal/entity"
	"RapidURL/internal/lib/random"
	repository "RapidURL/internal/repository/link"
)

type Usecase struct {
	s repository.Repository
}

func New(s repository.Repository) *Usecase {
	return &Usecase{
		s: s,
	}
}

const AliasLength = 5

func (u *Usecase) SaveLink(ctx context.Context, link entity.Link) (string, error) {
	if link.Alias == "" {
		link.Alias = random.NewRandomString(AliasLength)
	}

	return link.Alias, u.s.SaveLink(ctx, repository.DTO{
		Alias:  link.Alias,
		Url:    link.Url.String(),
		UserId: link.User.Id,
	})
}

func (u *Usecase) GetLink(ctx context.Context, alias string) (url.URL, error) {
	link, err := u.s.FindLinkByAlias(ctx, alias)

	if err != nil {
		return url.URL{}, err
	}

	return link.Url, err
}
