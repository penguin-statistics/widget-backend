package item

import (
	"time"

	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/utils"
)

// New creates a new Controller with its corresponding utils.Cache
func New() *Controller {
	logger := utils.NewLogger("ItemController")

	cache := utils.NewCache(utils.CacheConfig{
		Name:     "item",
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
func (c *Controller) Content() []*Item {
	return c.cache.Content().([]*Item)
}

// Status gives status.Status of current controller with the cache status
func (c *Controller) Status(_ string) *status.Status {
	return &status.Status{
		UpdatedAt: c.cache.Updated,
		FailCount: c.cache.FailCount,
		Length:    len(c.Content()),
	}
}

// Item returns the Item found with specified itemID
func (c *Controller) Item(itemID string) (result *Item) {
	for _, entry := range c.Content() {
		if entry.ItemID == itemID {
			return entry
		}
	}
	return nil
}
