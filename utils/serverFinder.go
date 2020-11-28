package utils

import "errors"

var (
	// ErrCacheNotFound describes a cache that cannot found with the server provided
	ErrCacheNotFound = errors.New("cache not found with the server provided")
)

// FindServerCache finds Cache that belongs to the specified server. Returns ErrCacheNotFound if not found.
func FindServerCache(caches []*Cache, server string) (*Cache, error) {
	for _, cache := range caches {
		if cache.Server == server {
			return cache, nil
		}
	}
	return nil, ErrCacheNotFound
}
