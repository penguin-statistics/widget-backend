package matrix

import (
	"time"

	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/utils"
)

// New creates a new Controller with its corresponding utils.Cache
func New() *Controller {
	logger := utils.NewLogger("MatrixController")

	var caches []*utils.Cache
	for _, server := range config.C.Upstream.Meta.Servers {
		caches = append(
			caches,
			utils.NewCache(utils.CacheConfig{
				Name:     "matrix",
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
func (c *Controller) ServerContent(server string) ([]*Matrix, *errors.Error) {
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

func filterStages(matrix []*Matrix, stageID string) (results []*Matrix) {
	for _, entry := range matrix {
		if entry.StageID == stageID {
			results = append(results, entry)
		}
	}
	return results
}

func filterItems(matrix []*Matrix, itemID string) (results []*Matrix) {
	for _, entry := range matrix {
		if entry.ItemID == itemID {
			results = append(results, entry)
		}
	}
	return results
}

// Query queries a set of matrix response using provided Query
func (c *Controller) Query(query *Query) (results []*Matrix, err *errors.Error) {
	data, sErr := c.ServerContent(query.Server)
	if sErr != nil {
		return nil, errors.New("FetchData", "failed to fetch matrix data from cache", errors.BlameServer)
	}
	unfiltered := len(data)

	if query.StageID != "" && query.ItemID != "" {
		results = filterStages(filterItems(data, query.ItemID), query.StageID)
	} else {
		if query.StageID != "" {
			results = filterStages(data, query.StageID)
		} else if query.ItemID != "" {
			results = filterItems(data, query.ItemID)
		}
	}

	if unfiltered == len(results) {
		return nil, errors.New("EmptyParams", "at least one of the query params shall be specified: stageId, itemId", errors.BlameUser)
	}

	if len(results) == 0 {
		return nil, errors.New("NotFound", "no records have been found with query params provided", errors.BlameUser)
	}

	return results, nil
}
