package utils

import "github.com/penguin-statistics/widget-backend/errors"

var (
	// ErrCacheNotFound describes a cache that cannot found with the server provided
	ErrCacheNotFound = errors.New("CacheNotFound", "cache not found with the server provided", errors.BlameUser)
)

// FindServerCache finds Cache that belongs to the specified server. Returns ErrCacheNotFound if not found.
func FindServerCache(caches []*Cache, server string) (*Cache, *errors.Error) {
	for _, cache := range caches {
		if cache.Server == server {
			return cache, nil
		}
	}
	return nil, ErrCacheNotFound
}
