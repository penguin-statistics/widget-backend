package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
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
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 8}))

	l.Debugln("`echo` has been initialized")

	l.Debugln("initializing controllers")

	controllers := meta.NewCollection()
	render := response.New(controllers, config.UILocation)

	matrixGroup := e.Group("/matrix/:server", func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			server := strings.ToUpper(c.Param("server"))
			if !config.ValidServer(server) {
				return c.HTMLBlob(render.Error(ErrInvalidServer))
			}

			c.Set("server", server)

			return handlerFunc(c)
		}
	})

	matrixGroup.GET("/stage/:stageId", func(c echo.Context) error {
		server := c.Get("server").(string)
		stageId := c.Param("stageId")
		records, err := controllers.Matrix.Stage(server, stageId)
		if err != nil {
			return c.HTMLBlob(render.Error(ErrFetchData))
		}
		resp, err := render.Marshal(records, server, &response.MatrixQuery{StageID: stageId})
		if err != nil {
			return c.HTMLBlob(render.Error(ErrCantMarshal))
		}
		return c.HTMLBlob(http.StatusOK, render.Response(resp))
	})

	matrixGroup.GET("/item/:itemId", func(c echo.Context) error {
		server := c.Get("server").(string)
		itemId := c.Param("itemId")
		records, err := controllers.Matrix.Item(server, itemId)
		if err != nil {
			return c.HTMLBlob(render.Error(ErrFetchData))
		}
		resp, err := render.Marshal(records, server, &response.MatrixQuery{ItemID: itemId})
		if err != nil {
			return c.HTMLBlob(render.Error(ErrCantMarshal))
		}
		return c.HTMLBlob(http.StatusOK, render.Response(resp))
	})

	matrixGroup.GET("/exact/:stageId/:itemId", func(c echo.Context) error {
		server := c.Get("server").(string)
		itemId := c.Param("itemId")
		stageId := c.Param("stageId")
		records, err := controllers.Matrix.Item(server, itemId)
		if err != nil {
			return c.HTMLBlob(render.Error(ErrFetchData))
		}

		var results []*matrix.Matrix
		for _, entry := range records {
			if entry.StageID == stageId {
				results = append(results, entry)
			}
		}
		resp, err := render.Marshal(results, server, &response.MatrixQuery{StageID: stageId, ItemID: itemId})
		if err != nil {
			return c.HTMLBlob(render.Error(ErrCantMarshal))
		}
		return c.HTMLBlob(http.StatusOK, render.Response(resp))
	})

	e.GET("/item/:itemId", func(c echo.Context) error {
		data, err := controllers.Item.Item(c.Param("itemId"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, data)
	})

	e.Static("/", config.UILocation)

	l.Fatalln(e.Start(":8000"))
}
