package stage

import (
	"time"

	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/utils"
)

// New creates a new Controller with its corresponding utils.Cache
func New() *Controller {
	logger := utils.NewLogger("StageController")

	var caches []*utils.Cache
	for _, server := range config.C.Upstream.Meta.Servers {
		caches = append(
			caches,
			utils.NewCache(utils.CacheConfig{
				Name:     "stage",
				Server:   server,
				Interval: time.Minute * 30,
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
func (c *Controller) ServerContent(server string) ([]*Stage, *errors.Error) {
	cache, err := c.Server(server)
	if err != nil {
		return nil, err
	}
	return cache.Content().([]*Stage), nil
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
		Length:    len(inst.Content().([]*Stage)),
	}
}

// Stage returns the Stage found with specified stageID in server
func (c *Controller) Stage(server string, stageID string) (result *Stage) {
	data, err := c.ServerContent(server)
	if err != nil {
		return nil
	}

	for _, entry := range data {
		if entry.StageID == stageID {
			return entry
		}
	}
	return nil
}
