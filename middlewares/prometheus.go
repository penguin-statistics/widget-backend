package middlewares

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"strings"
)

func prometheusSkipper(c echo.Context) bool {
	path := c.Request().URL.Path
	return strings.HasPrefix(path, "_")
}

func NewPrometheus(e *echo.Echo) {
	prom := prometheus.NewPrometheus("widget", prometheusSkipper)
	prom.Use(e)
}
