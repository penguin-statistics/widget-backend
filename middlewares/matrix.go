package middlewares

import (
	"github.com/biter777/countries"
	"github.com/labstack/echo/v4"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/siteStats"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/response"
	"github.com/penguin-statistics/widget-backend/utils"
)

// SiteStatsQuery is an echo.MiddlewareFunc that verifies siteStats.Query parameters and inject to echo.Context
func SiteStatsQuery(render *response.Assembler) func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			server := c.Param("server")
			if !utils.ValidServer(server) {
				return c.HTMLBlob(render.HTMLError(errors.ErrInvalidServer))
			}

			// populate matrix query preferences
			c.Set("query", &siteStats.Query{
				Server:  server,
			})

			return handlerFunc(c)
		}
	}
}

// MatrixQuery is an echo.MiddlewareFunc that verifies basic prerequisites that the request shall comply with and injects the formatted query request into echo.Context
func MatrixQuery(render *response.Assembler) func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			server := c.Param("server")
			if !utils.ValidServer(server) {
				return c.HTMLBlob(render.HTMLError(errors.ErrInvalidServer))
			}

			// populate matrix query preferences
			c.Set("query", &matrix.Query{
				Server:  server,
				StageID: c.Param("stageId"),
				ItemID:  c.Param("itemId"),
			})

			return handlerFunc(c)
		}
	}
}

// RequestMetadata is an echo.MiddlewareFunc that populates response.RequestMetadata and injects it into echo.Context
func RequestMetadata() func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// populate response.RequestMetadata fields
			// by default we don't provide mirror value
			mirror := ""
			ipCountry := c.Request().Header.Get("CF-IPCountry")
			// if user is in China Mainland, a "cn" mirror is preferred rather the "io" mirror
			if ipCountry == countries.CN.Alpha2() {
				mirror = "cn"
			} else if ipCountry != "" {
				mirror = "io"
			}

			c.Set("meta", &response.RequestMetadata{
				Mirror: mirror,
			})

			return handlerFunc(c)
		}
	}
}
