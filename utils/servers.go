package utils

import (
	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/errors"
)

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

// ValidServer validates if server counts as a valid server
func ValidServer(server string) bool {
	for _, s := range config.C.Upstream.Meta.Servers {
		if s == server {
			return true
		}
	}
	return false
}