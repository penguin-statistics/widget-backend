package siteStats

import (
	"github.com/penguin-statistics/widget-backend/utils"
	"io"
	"net/http"
	"net/url"
	"sort"
)

var logger = utils.NewLogger("SiteStatsController:utils")

func unmarshalResponse(reader io.Reader) (result []StageTime, err error) {
	var response RemoteResponse
	err = utils.UnmarshalFromReader(reader, &response)
	if err != nil {
		logger.Errorln("failed to unmarshal siteStats response from reader", err)
		return nil, err
	}

	stages := response.TotalStageTimesRecent24H

	sort.SliceStable(stages, func(i, j int) bool {
		return stages[i].RecentTimes > stages[j].RecentTimes
	})

	return stages, err
}

func createUpdater(server string) utils.Updater {
	return func(cache *utils.Cache) (interface{}, error) {
		params := url.Values{}
		params.Set("server", server)

		req, err := utils.NewPenguinRequest("stats", &params)
		if err != nil {
			panic(err)
		}

		cache.Logger.Traceln("assembled new request with url", req.URL)

		var resp *http.Response
		err = utils.NewRetriedOperation(func() (err error) {
			resp, err = cache.Client.Do(req)
			return err
		})

		if err != nil {
			cache.Logger.Errorln("failed to fetch external data after multiple retries.", err)
			return nil, err
		}

		return unmarshalResponse(resp.Body)
	}
}
