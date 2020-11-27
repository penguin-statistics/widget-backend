package item

import (
	"github.com/penguin-statistics/partial-matrix/utils"
	"time"
)

func New() *Controller {
	logger := utils.NewLogger("ItemController")

	cache := utils.NewCache(utils.CacheConfig{
		Name:     "Item",
		Server:   "",
		Interval: time.Minute * 30,
		Updater:  updater,
	})

	return &Controller{
		cache:  cache,
		logger: logger,
	}
}

func (c *Controller) Content() []*Item {
	return c.cache.Content().([]*Item)
}

func (c *Controller) Item(itemId string) (result *Item, err error) {
	for _, entry := range c.Content() {
		if entry.ItemID == itemId {
			return entry, nil
		}
	}
	return nil, nil
}
