package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/meta"
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/middlewares"
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
	e.Use(middlewares.RequestMetadata())

	l.Debugln("`echo` has been initialized")

	l.Debugln("initializing controllers...")

	controllers := meta.NewCollection()

	l.Debugln("controllers initialized. initializing render...")
	render := response.New(controllers, config.UILocation)

	l.Debugln("render initialized. registering handlers...")

	// HTML Rendered Response
	{
		rendered := e.Group("/result/:server", middlewares.MatrixQuery(render))

		renderedHandler := func(c echo.Context) error {
			query := c.Get("query").(*matrix.Query)

			records, err := controllers.Matrix.Query(query)
			if err != nil {
				return c.HTMLBlob(render.HTMLError(err))
			}
			return render.HTMLResponse(c, render.Marshal(records, query))
		}

		rendered.GET("/stage/:stageId", renderedHandler)
		rendered.GET("/item/:itemId", renderedHandler)
		rendered.GET("/exact/:stageId/:itemId", renderedHandler)
	}

	// API Response
	{
		api := e.Group("/api/result/:server", middlewares.MatrixQuery(render))

		apiHandler := func(c echo.Context) error {
			query := c.Get("query").(*matrix.Query)

			records, err := controllers.Matrix.Query(query)
			if err != nil {
				return c.JSON(render.JSONError(err))
			}
			return render.JSONResponse(c, render.Marshal(records, query))
		}

		api.GET("/stage/:stageId", apiHandler)
		api.GET("/item/:itemId", apiHandler)
		api.GET("/exact/:stageId/:itemId", apiHandler)
	}

	e.GET("/_health", func(c echo.Context) error {
		var statusInd int
		statuses := map[string]map[string]*status.Status{}
		for _, server := range config.Server {
			serversStatus := controllers.Statuses(server)
			for _, serverStatus := range serversStatus {
				statusInd += serverStatus.FailCount
			}
			statuses[server] = serversStatus
		}

		httpStatus := http.StatusOK
		// 4: cache type instances
		if statusInd >= 4 {
			httpStatus = http.StatusServiceUnavailable
		}
		return c.JSON(httpStatus, struct {
			Status int `json:"status"`
			CacheStatuses map[string]map[string]*status.Status `json:"caches"`
			System SystemMetrics `json:"system"`
		} {
			Status: statusInd,
			CacheStatuses: statuses,
			System: newSystemMetrics(),
		})
	})

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
		return c.HTMLBlob(render.HTMLError(errors.New("PageNotFound", "unrecognized resource path", errors.BlameUser)))
	})

	//l.Traceln(spew.Sdump(e.Routes()))

	l.Debugln("handlers registered. starting http server...")

	l.Fatalln(e.Start(":8010"))
}
