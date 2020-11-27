package stage

import (
	"github.com/penguin-statistics/partial-matrix/utils"
	"io"
	"net/http"
)

var logger = utils.NewLogger("StageController:utils")

func unmarshalResponse(reader io.Reader) (result []*Stage, err error) {
	err = utils.UnmarshalFromReader(reader, &result)
	if err != nil {
		logger.Errorln("failed to unmarshal matrix response from reader", err)
		return nil, err
	}
	return result, err
}

func updater(cache *utils.Cache) (interface{}, error) {
	req, err := utils.NewPenguinRequest("stages", "")
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

	return unmarshalResponse(resp.Body)
}
