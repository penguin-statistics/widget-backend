package zone

import (
	"github.com/penguin-statistics/partial-matrix/utils"
	"time"
)

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

func (c *Controller) Content() []*Zone {
	return c.cache.Content().([]*Zone)
}

func (c *Controller) Zone(zoneId string) (result *Zone, err error) {
	for _, entry := range c.Content() {
		if entry.ZoneID == zoneId {
			return entry, nil
		}
	}
	return nil, nil
}
