package middlewares

import (
	"github.com/labstack/echo"
)

// CacheType describes the type of cache shall be used, which its value is the content of the `Cache-Control` header that is being populated
type CacheType string

const (
	// CacheTypeStatic describes the resource should never be changed and therefore is safe for long-term caching
	CacheTypeStatic CacheType = "public, max-age=31356000"
	// CacheTypeDynamic describes the resource as if it is a valuable-to-cache dynamic resource
	CacheTypeDynamic CacheType = "public, max-age=300, must-revalidate"
	// CacheTypeNoCache describes the resource should always not be cached
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
