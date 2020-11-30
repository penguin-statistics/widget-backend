package stage

import (
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/utils"
	"time"
)

// New creates a new Controller with its corresponding utils.Cache
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

// Content returns the content from its cache
func (c *Controller) Content() []*Stage {
	return c.cache.Content().([]*Stage)
}

// Status gives status.Status of current controller with the cache status
func (c *Controller) Status(_ string) *status.Status {
	return &status.Status{
		UpdatedAt: c.cache.Updated,
		FailCount: c.cache.FailCount,
		Length:    len(c.Content()),
	}
}

// Stage returns the Stage found with specified stageID
func (c *Controller) Stage(stageID string) (result *Stage) {
	for _, entry := range c.Content() {
		if entry.StageID == stageID {
			return entry
		}
	}
	return nil
}
