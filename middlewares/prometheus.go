package middlewares

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func prometheusSkipper(c echo.Context) bool {
	if c.Response().Status == http.StatusNotFound {
		return true
	}
	path := c.Request().URL.Path
	return strings.HasPrefix(path, "/_")
}

func NewPrometheus(e *echo.Echo) {
	prom := prometheus.NewPrometheus("echo", prometheusSkipper)
	prom.Use(e)
}
