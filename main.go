package main

import (
	"embed"
	"net/http"
	"path"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/meta"
	"github.com/penguin-statistics/widget-backend/controller/siteStats"
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/middlewares"
	"github.com/penguin-statistics/widget-backend/response"
	"github.com/penguin-statistics/widget-backend/utils"
)

//go:embed misc/teapot
var teapot embed.FS

var l = utils.NewLogger("main")

func init() {
	l.Infoln("launching...")
}

func main() {
	// echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
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
	middlewares.NewPrometheus(e)

	l.Debugln("config loaded as", config.C)

	l.Debugln("`echo` has been initialized")

	l.Debugln("initializing controllers...")

	controllers := meta.NewCollection()

	l.Debugln("controllers initialized. initializing render...")
	render := response.New(controllers, config.C.Static.Widget.Root)

	l.Debugln("render initialized. registering handlers...")

	// HTML Rendered Response (Matrix)
	{
		rendered := e.Group("/result/:server", middlewares.MatrixQuery(render), middlewares.PopulateCacheHeader(middlewares.CacheTypeDynamic))

		renderedHandler := func(c echo.Context) error {
			query := c.Get("query").(*matrix.Query)

			records, err := controllers.Matrix.Query(query)
			if err != nil {
				return c.HTMLBlob(render.HTMLError(err))
			}
			return render.HTMLResponse(c, render.MarshalMatrix(records, query))
		}

		rendered.GET("/stage/:stageId", renderedHandler)
		rendered.GET("/item/:itemId", renderedHandler)
		rendered.GET("/exact/:stageId/:itemId", renderedHandler)
	}

	// API Response (Matrix)
	{
		api := e.Group("/api/result/:server", middlewares.MatrixQuery(render), middlewares.PopulateCacheHeader(middlewares.CacheTypeDynamic))

		apiHandler := func(c echo.Context) error {
			query := c.Get("query").(*matrix.Query)

			records, err := controllers.Matrix.Query(query)
			if err != nil {
				return c.JSON(render.JSONError(err))
			}
			return render.JSONMatrixResponse(c, render.MarshalMatrix(records, query))
		}

		api.GET("/stage/:stageId", apiHandler)
		api.GET("/item/:itemId", apiHandler)
		api.GET("/exact/:stageId/:itemId", apiHandler)
	}

	// API Response (SiteStats)
	{
		statsApi := e.Group("/api/stats/:server", middlewares.SiteStatsQuery(render), middlewares.PopulateCacheHeader(middlewares.CacheTypeDynamic))
		statsApi.GET("", func(c echo.Context) error {
			query := c.Get("query").(*siteStats.Query)

			records, err := controllers.SiteStats.Query(query)
			if err != nil {
				return c.JSON(render.JSONError(err))
			}
			return c.JSON(http.StatusOK, render.MarshalSiteStats(records, query))
		})
	}

	e.GET("/_health", func(c echo.Context) error {
		var statusInd int
		statuses := map[string]map[string]*status.Status{}
		for _, server := range config.C.Upstream.Meta.Servers {
			serversStatus := controllers.Statuses(server)
			for _, serverStatus := range serversStatus {
				statusInd += serverStatus.FailCount
			}
			statuses[server] = serversStatus
		}

		httpStatus := http.StatusOK
		if statusInd >= len(statuses) {
			httpStatus = http.StatusServiceUnavailable
		}
		return c.JSON(httpStatus, struct {
			Status        int                                  `json:"status"`
			CacheStatuses map[string]map[string]*status.Status `json:"caches"`
		}{
			Status:        statusInd,
			CacheStatuses: statuses,
		})
	}, middlewares.PopulateCacheHeader(middlewares.CacheTypeNoCache))

	{
		static := e.Group("/", middlewares.PopulateCacheHeader(middlewares.CacheTypeStatic))

		// widget static files
		static.Static(config.C.Static.Widget.Endpoint, path.Join(config.C.Static.Widget.Root, "_widget"))

		// docs static files
		static.Static(config.C.Static.Docs.Endpoint, path.Join(config.C.Static.Docs.Root, "_docs"))
	}

	// docs page
	e.File("/", path.Join(config.C.Static.Docs.Root, "index.html"))

	// favicon which directs to /favicon.ico at widget frontend
	e.File("/favicon.ico", path.Join(config.C.Static.Widget.Root, "favicon.ico"))

	// very important to implement. without this the server won't behave normally. (of course not) :D
	e.Any("/_teapot", func(c echo.Context) error {
		b, err := teapot.ReadFile("misc/teapot")
		if err != nil {
			return c.JSON(render.JSONError(errors.New("TeapotNotFound", "cannot get teapot due to server failure :(", errors.BlameServer)))
		}
		return c.Blob(http.StatusTeapot, echo.MIMETextPlainCharsetUTF8, b)
	}, middlewares.PopulateCacheHeader(middlewares.CacheTypeDynamic))

	// display customized error handler page
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if !c.Response().Committed {
			c.Response().Header().Set("Cache-Control", "no-store")
			c.Response().Header().Set("X-Robots-Tag", "noindex")

			if err == echo.ErrNotFound && c.Request().Method != http.MethodHead {
				rErr := c.HTMLBlob(render.HTMLError(errors.New("PageNotFound", "unrecognized resource path", http.StatusNotFound)))
				if rErr != nil {
					l.Errorln("failed to return custom handler page", err)
				}
			} else {
				e.DefaultHTTPErrorHandler(err, c)
			}
		}
	}

	l.Debugln("handlers registered. starting http server...")

	l.Fatalln(e.Start(config.C.Server.Listen))
}
