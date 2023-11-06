package link

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"RapidURL/internal/entity"
	"RapidURL/internal/lib/logger/sl"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/prometheus/client_golang/prometheus"
)

var _ Repository = &CachedRepository{}

type CachedRepository struct {
	repository       Repository
	cache            CacheRepository
	log              *slog.Logger
	totalCacheMisses prometheus.Counter
}

// NewCachedRepository receive a link repository and cache repository
// and return a composition of it
func NewCachedRepository(repository Repository, cache CacheRepository, log *slog.Logger, totalCacheMisses prometheus.Counter) *CachedRepository {
	return &CachedRepository{
		repository:       repository,
		cache:            cache,
		log:              log,
		totalCacheMisses: totalCacheMisses,
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
			c.totalCacheMisses.Add(1)
		} else {
			c.log.Error("error while trying find link in cache", sl.Err(err))
		}

		link, err = c.repository.FindLinkByAlias(ctx, alias)
		c.cache.SaveLink(ctx, FromEntity(link))
		return link, err
	}

	return link, nil
}
