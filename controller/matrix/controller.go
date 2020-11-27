package matrix

import (
	"github.com/penguin-statistics/partial-matrix/config"
	"github.com/penguin-statistics/partial-matrix/utils"
	"time"
)

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

func (c *Controller) Server(server string) ([]*Matrix, error) {
	cache, err := utils.FindServerCache(c.caches, server)
	if err != nil {
		return nil, err
	}
	return cache.Content().([]*Matrix), nil
}

func (c *Controller) Stage(server, stageId string) (results []*Matrix, err error) {
	data, err := c.Server(server)
	if err != nil {
		return nil, err
	}

	for _, entry := range data {
		if entry.StageID == stageId {
			results = append(results, entry)
		}
	}
	return results, nil
}

func (c *Controller) Item(server, itemId string) (results []*Matrix, err error) {
	data, err := c.Server(server)
	if err != nil {
		return nil, err
	}

	for _, entry := range data {
		if entry.ItemID == itemId {
			results = append(results, entry)
		}
	}
	return results, nil
}
