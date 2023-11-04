package link

import (
	"net/url"

	"RapidURL/internal/entity"
	"RapidURL/internal/lib/random"
)

type Storage interface {
	SaveLink(link *entity.Link) error
	FindLinkByAlias(alias string) (*entity.Link, error)
}

type Usecase struct {
	s Storage
}

func New(s Storage) *Usecase {
	return &Usecase{
		s: s,
	}
}

const AliasLength = 5

func (u *Usecase) SaveLink(link SaveLinkDTO) (string, error) {
	if link.Alias == "" {
		link.Alias = random.NewRandomString(AliasLength)
	}

	return link.Alias, u.s.SaveLink(&entity.Link{
		Alias:  link.Alias,
		Url:    &link.Url,
		UserId: link.UserId,
	})
}

func (u *Usecase) GetLink(dto GetLinkDTO) (url.URL, error) {
	link, err := u.s.FindLinkByAlias(dto.Alias)

	if err != nil {
		return url.URL{}, err
	}

	return *link.Url, err
}
