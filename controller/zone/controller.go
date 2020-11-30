package zone

import (
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/utils"
	"time"
)

// New creates a new Controller with its corresponding utils.Cache
func New() *Controller {
	logger := utils.NewLogger("ZoneController")

	cache := utils.NewCache(utils.CacheConfig{
		Name:     "Zone",
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
func (c *Controller) Content() []*Zone {
	return c.cache.Content().([]*Zone)
}

// Status gives status.Status of current controller with the cache status
func (c *Controller) Status(_ string) *status.Status {
	return &status.Status{
		UpdatedAt: c.cache.Updated,
		FailCount: c.cache.FailCount,
		Length:    len(c.Content()),
	}
}

// Zone returns the Zone found with specified zoneID
func (c *Controller) Zone(zoneID string) (result *Zone) {
	for _, entry := range c.Content() {
		if entry.ZoneID == zoneID {
			return entry
		}
	}
	return nil
}
