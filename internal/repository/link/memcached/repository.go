package memcached

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"RapidURL/internal/entity"
	"RapidURL/internal/repository/link"
	"github.com/bradfitz/gomemcache/memcache"
)

const expiration = 60 * 60 * 24
const timeout = 500 * time.Millisecond

type Repository struct {
	client *memcache.Client
}

func New(client *memcache.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) SaveLink(ctx context.Context, dto link.DTO) error {
	const op = "repository.link.memcached.SaveLink"

	r.setTimeOut(ctx)
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

	r.setTimeOut(ctx)
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

func (r *Repository) setTimeOut(ctx context.Context) {
	deadline, _ := ctx.Deadline()
	if time.Until(deadline) > 0 {
		r.client.Timeout = time.Until(deadline)
	} else {
		r.client.Timeout = 0
	}
}
