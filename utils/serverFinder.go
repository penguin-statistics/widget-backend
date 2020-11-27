package utils

import "errors"

func FindServerCache(caches []*Cache, server string) (*Cache, error) {
	for _, cache := range caches {
		if cache.Server == server {
			return cache, nil
		}
	}
	return nil, errors.New("cache not found with the server provided")
}
