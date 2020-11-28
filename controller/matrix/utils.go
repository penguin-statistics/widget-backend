package matrix

import (
	"github.com/penguin-statistics/widget-backend/utils"
	"io"
	"net/http"
)

var logger = utils.NewLogger("MatrixController:utils")

func unmarshalResponse(reader io.Reader) (result []*Matrix, err error) {
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

func createUpdater(server string) utils.Updater {
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

		return unmarshalResponse(resp.Body)
	}
}
