package middlewares

import (
	"github.com/biter777/countries"
	"github.com/labstack/echo"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/response"
	"github.com/penguin-statistics/widget-backend/utils"
)

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
			// by default mirror is "io"
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
