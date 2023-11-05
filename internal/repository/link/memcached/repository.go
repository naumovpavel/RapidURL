package memcached

import (
	"context"
	"encoding/json"
	"fmt"

	"RapidURL/internal/entity"
	"RapidURL/internal/repository/link"
	"github.com/bradfitz/gomemcache/memcache"
)

const expiration = 60 * 60 * 24

type Repository struct {
	client *memcache.Client
}

func New(client *memcache.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) SaveLink(ctx context.Context, dto link.DTO) error {
	const op = "repository.link.memcached.SaveLink"
	json, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return r.client.Add(&memcache.Item{
		Key:        dto.Alias,
		Value:      json,
		Expiration: expiration,
	})
}
func (r *Repository) FindLinkByAlias(ctx context.Context, alias string) (entity.Link, error) {
	const op = "repository.link.memcached.FindLinkByAlias"
	item, err := r.client.Get(alias)
	if err != nil {
		return entity.Link{}, fmt.Errorf("%s: %w", op, err)
	}

	var dto link.DTO
	err = json.Unmarshal(item.Value, &dto)
	if err != nil {
		return entity.Link{}, fmt.Errorf("%s: %w", op, err)
	}

	return link.ToEntity(dto)
}
