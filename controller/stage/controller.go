package stage

import (
	"github.com/penguin-statistics/partial-matrix/utils"
	"time"
)

func New() *Controller {
	logger := utils.NewLogger("StageController")

	cache := utils.NewCache(utils.CacheConfig{
		Name:     "Stage",
		Server:   "",
		Interval: time.Minute * 30,
		Updater:  updater,
	})

	return &Controller{
		cache:  cache,
		logger: logger,
	}
}

func (c *Controller) Content() []*Stage {
	return c.cache.Content().([]*Stage)
}

func (c *Controller) Stage(stageId string) (result *Stage, err error) {
	for _, entry := range c.Content() {
		if entry.StageID == stageId {
			return entry, nil
		}
	}
	return nil, nil
}
