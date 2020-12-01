package stage

import (
	"github.com/penguin-statistics/widget-backend/utils"
	"io"
	"net/http"
	"net/url"
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

func createUpdater(server string) utils.Updater {
	return func(cache *utils.Cache) (interface{}, error) {
		params := url.Values{}
		params.Set("server", server)

		req, err := utils.NewPenguinRequest("stages", &params)
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

//func updater(cache *utils.Cache) (interface{}, error) {
//	req, err := utils.NewPenguinRequest("stages", &url.Values{})
//	if err != nil {
//		panic(err)
//	}
//
//	cache.Logger.Traceln("assembled new request with url", req.URL)
//
//	var resp *http.Response
//	err = utils.NewRetriedOperation(func() (err error) {
//		resp, err = cache.Client.Do(req)
//		return err
//	})
//
//	if err != nil {
//		cache.Logger.Errorln("failed to fetch external data after multiple retries.", err)
//		return nil, err
//	}
//
//	return unmarshalResponse(resp.Body)
//}
