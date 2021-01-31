package siteStats

import (
	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/utils"
	"time"
)

const (
	MaxSiteStatsAmount = 3
)

// New creates a new Controller with its corresponding utils.Cache
func New() *Controller {
	logger := utils.NewLogger("SiteStatsController")

	var caches []*utils.Cache
	for _, server := range config.C.Upstream.Meta.Servers {
		caches = append(
			caches,
			utils.NewCache(utils.CacheConfig{
				Name:     "SiteStat",
				Server:   server,
				Interval: time.Minute * 5,
				Updater:  createUpdater(server),
			}),
		)
		logger.Debugln("cache created for server", server)
	}

	return &Controller{
		caches: caches,
		logger: logger,
	}
}

// Server gives utils.Cache of server
func (c *Controller) Server(server string) (*utils.Cache, *errors.Error) {
	cache, err := utils.FindServerCache(c.caches, server)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// ServerContent gives matrices from utils.Cache of server
func (c *Controller) ServerContent(server string) ([]StageTime, *errors.Error) {
	cache, err := c.Server(server)
	if err != nil {
		return nil, err
	}
	return cache.Content().([]StageTime), nil
}

// Status gives status.Status of current controller with the status of utils.Cache for server
func (c *Controller) Status(server string) *status.Status {
	inst, err := c.Server(server)
	if err != nil {
		return &status.Status{}
	}
	return &status.Status{
		UpdatedAt: inst.Updated,
		FailCount: inst.FailCount,
		Length:    len(inst.Content().([]StageTime)),
	}
}

// Query queries a set of siteStats response using provided Query
func (c *Controller) Query(query *Query) (results []StageTime, err *errors.Error) {
	data, sErr := c.ServerContent(query.Server)
	if sErr != nil {
		return nil, errors.New("FetchData", "failed to fetch siteStats data from cache", errors.BlameServer)
	}

	results = data[:MaxSiteStatsAmount]

	if len(results) == 0 {
		return nil, errors.New("NotFound", "no records have been found with query params provided", errors.BlameUser)
	}

	return results, nil
}
