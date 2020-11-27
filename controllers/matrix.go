package controllers

import (
	"github.com/penguin-statistics/partial-matrix/config"
	"github.com/penguin-statistics/partial-matrix/utils"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

var logger = utils.NewLogger("MatrixController:utils")

type MatrixController struct {
	caches []*utils.Cache
	logger *logrus.Entry
}

type Matrix struct {
	StageID  string `json:"stageId,omitempty"`
	ItemID   string `json:"itemId,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
	Times    int    `json:"times,omitempty"`
	Start    *int   `json:"start,omitempty"`
	End      *int   `json:"end,omitempty"`
}

func unmarshalMatrixResponse(reader io.Reader) (result []*Matrix, err error) {
	var response struct {
		Matrix []*Matrix `json:"matrix"`
	}
	err = utils.UnmarshalFromReader(reader, &response)
	if err != nil {
		logger.Errorln("failed to unmarshal matrix response from reader", err)
		return nil, err
	}
	return response.Matrix, err
}

func createMatrixUpdater(server string) utils.Updater {
	return func(cache *utils.Cache) (interface{}, error) {
		req, err := utils.NewPenguinRequest("result/matrix", server)
		if err != nil {
			panic(err)
		}

		cache.Logger.Debugln("assembled new request with url", req.URL)

		var resp *http.Response
		err = utils.NewRetriedOperation(func() (err error) {
			resp, err = cache.Client.Do(req)
			return err
		})

		if err != nil {
			cache.Logger.Errorln("failed to fetch external data after multiple retries.", err)
			return nil, err
		}

		return unmarshalMatrixResponse(resp.Body)
	}
}

func NewMatrixController() *MatrixController {
	logger := utils.NewLogger("MatrixController")

	var caches []*utils.Cache
	for _, server := range config.Server {
		caches = append(
			caches,
			utils.NewCache(utils.CacheConfig{
				Name: "Matrix",
				Server: server,
				Interval: time.Minute * 5,
				Updater: createMatrixUpdater(server),
			}),
		)
		logger.Debugln("cache created for server", server)
	}

	return &MatrixController{
		caches: caches,
		logger: logger,
	}
}

func (c *MatrixController) Server(server string) ([]*Matrix, error) {
	cache, err := utils.FindServerCache(c.caches, server)
	if err != nil {
		return nil, err
	}
	return cache.Content().([]*Matrix), nil
}

func (c *MatrixController) Stage(server, stageId string) (results []*Matrix, err error) {
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

func (c *MatrixController) Item(server, itemId string) (results []*Matrix, err error) {
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
