package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/meta"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/response"
	"github.com/penguin-statistics/widget-backend/utils"
	"net/http"
	"path"
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
		AllowCredentials: true,
		ExposeHeaders: []string{
			"Last-Modified",
		},
		MaxAge: int((time.Hour * 24 * 365).Seconds()),
	}))

	l.Debugln("`echo` has been initialized")

	l.Debugln("initializing controllers...")

	controllers := meta.NewCollection()

	l.Debugln("controllers initialized. initializing render...")
	render := response.New(controllers, config.UILocation)

	l.Debugln("render initialized. registering handlers...")

	matrixGroup := e.Group("/result/:server", func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			server := c.Param("server")
			if !config.ValidServer(server) {
				return c.HTMLBlob(render.Error(ErrInvalidServer))
			}

			c.Set("query", &matrix.Query{
				Server:  server,
				StageID: c.Param("stageId"),
				ItemID:  c.Param("itemId"),
			})

			return handlerFunc(c)
		}
	})

	matrixHandler := func(c echo.Context) error {
		query := c.Get("query").(*matrix.Query)

		records, err := controllers.Matrix.Query(query)
		if err != nil {
			return c.HTMLBlob(render.Error(err))
		}
		return render.Response(c, render.Marshal(records, query))
	}

	matrixGroup.GET("/stage/:stageId", matrixHandler)
	matrixGroup.GET("/item/:itemId", matrixHandler)
	matrixGroup.GET("/exact/:stageId/:itemId", matrixHandler)

	// widget static files
	e.Static("/_widget", path.Join(config.UILocation, "_widget"))

	// docs static files
	e.Static("/_docs", path.Join(config.DocLocation, "_docs"))

	// docs page
	e.File("/", path.Join(config.DocLocation, "index.html"))

	// favicon which directs to /favicon.ico at widget frontend
	e.File("/favicon.ico", path.Join(config.UILocation, "favicon.ico"))

	// match all other routes as 404 and display custom rendered error page
	e.GET("*", func(c echo.Context) error {
		return c.HTMLBlob(render.Error(errors.New("PageNotFound", "unrecognized resource path", errors.BlameUser)))
	})

	//l.Traceln(spew.Sdump(e.Routes()))

	l.Debugln("handlers registered. starting http server...")

	l.Fatalln(e.Start(":8010"))
}
