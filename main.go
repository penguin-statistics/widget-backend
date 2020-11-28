package main

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/controller/meta"
	"github.com/penguin-statistics/widget-backend/response"
	"github.com/penguin-statistics/widget-backend/utils"
	"net/http"
	"strings"
	"time"
)

var l = utils.NewLogger("main")

func init() {
	l.Infoln("launching...")
}

func main() {
	// echo
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"*",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Penguin-Source",
		},
		AllowCredentials: true,
		ExposeHeaders: []string{
			"Last-Modified",
			"Penguin-Record-Size",

			"RateLimit-Limit",
			"RateLimit-Remaining",
			"RateLimit-Reset",
		},
		MaxAge: int((time.Hour * 24 * 365).Seconds()),
	}))

	l.Debugln("`echo` has been initialized")

	l.Debugln("initializing controllers")

	controllers := meta.NewCollection()
	manager := response.New(controllers)

	matrixGroup := e.Group("/matrix/:server", func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			server := strings.ToUpper(c.Param("server"))
			if !config.ValidServer(server) {
				return echo.NewHTTPError(http.StatusBadRequest, errors.New("malformed parameter `server` provided"))
			}

			c.Set("server", server)

			return handlerFunc(c)
		}
	})

	matrixGroup.GET("/stage/:stageId", func(c echo.Context) error {
		server := c.Get("server").(string)
		records, err := controllers.Matrix.Stage(server, c.Param("stageId"))
		if err != nil {
			return err
		}
		resp, err := manager.Marshal(records, server, "stage")
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, resp)
	})

	matrixGroup.GET("/item/:itemId", func(c echo.Context) error {
		server := c.Get("server").(string)
		records, err := controllers.Matrix.Item(server, c.Param("itemId"))
		if err != nil {
			return err
		}
		resp, err := manager.Marshal(records, server, "item")
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, resp)
	})

	e.GET("/item/:itemId", func(c echo.Context) error {
		data, err := controllers.Item.Item(c.Param("itemId"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, data)
	})

	l.Fatalln(e.Start(":8000"))
}
