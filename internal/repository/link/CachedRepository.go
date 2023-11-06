package link

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"RapidURL/internal/entity"
	"RapidURL/internal/lib/logger/sl"
	"github.com/bradfitz/gomemcache/memcache"
)

var _ Repository = &CachedRepository{}

type CachedRepository struct {
	repository Repository
	cache      CacheRepository
	log        *slog.Logger
}

// NewCachedRepository receive a link repository and cache repository
// and return a composition of it
func NewCachedRepository(repository Repository, cache CacheRepository, log *slog.Logger) *CachedRepository {
	return &CachedRepository{
		repository: repository,
		cache:      cache,
		log:        log,
	}
}

func (c *CachedRepository) SaveLink(ctx context.Context, link DTO) error {
	const op = "repository.link.memcached.SaveLink"

	err := c.repository.SaveLink(ctx, link)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = c.cache.SaveLink(ctx, link)
	if err != nil {
		c.log.Error("can't save link in cache", sl.Err(fmt.Errorf("%s: %w", op, err)))
	}

	return nil
}

func (c *CachedRepository) FindLinkByAlias(ctx context.Context, alias string) (entity.Link, error) {
	const op = "repository.link.memcached.SaveLink"

	link, err := c.cache.FindLinkByAlias(ctx, alias)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			c.log.Warn("link cache miss", sl.Err(err))
		} else {
			c.log.Error("error while trying find link in cache", sl.Err(err))
		}

		return c.repository.FindLinkByAlias(ctx, alias)
	}

	return link, nil
}