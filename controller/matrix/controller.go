package matrix

import (
	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/utils"
	"time"
)

// New creates a new Controller with its corresponding utils.Cache
func New() *Controller {
	logger := utils.NewLogger("MatrixController")

	var caches []*utils.Cache
	for _, server := range config.Server {
		caches = append(
			caches,
			utils.NewCache(utils.CacheConfig{
				Name:     "Matrix",
				Server:   server,
				Interval: time.Minute * 1,
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
func (c *Controller) Server(server string) (*utils.Cache, error) {
	cache, err := utils.FindServerCache(c.caches, server)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// ServerContent gives matrices from utils.Cache of server
func (c *Controller) ServerContent(server string) ([]*Matrix, error) {
	cache, err := c.Server(server)
	if err != nil {
		return nil, err
	}
	return cache.Content().([]*Matrix), nil
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
		Length:    len(inst.Content().([]*Matrix)),
	}
}

// Stage returns the matrices found in utils.Cache for server with specified stageID
func (c *Controller) Stage(server, stageID string) (results []*Matrix, err error) {
	data, err := c.ServerContent(server)
	if err != nil {
		return nil, err
	}

	for _, entry := range data {
		if entry.StageID == stageID {
			results = append(results, entry)
		}
	}
	return results, nil
}

// Item returns the matrices found in utils.Cache for server with specified itemId
func (c *Controller) Item(server, itemID string) (results []*Matrix, err error) {
	data, err := c.ServerContent(server)
	if err != nil {
		return nil, err
	}

	for _, entry := range data {
		if entry.ItemID == itemID {
			results = append(results, entry)
		}
	}
	return results, nil
}
