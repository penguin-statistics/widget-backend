package middlewares

import (
	"github.com/labstack/echo"
)

type CacheType string

const (
	CacheTypeStatic CacheType = "public, max-age=31356000"
	CacheTypeDynamic CacheType = "public, max-age=300, must-revalidate"
	CacheTypeNoCache CacheType = "no-store"
)

// PopulateCacheHeader populates the `Cache-Control`
func PopulateCacheHeader(cacheType CacheType) func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", string(cacheType))

			return handlerFunc(c)
		}
	}
}
