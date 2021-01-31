package stage

import (
	"encoding/json"
	"github.com/penguin-statistics/widget-backend/models"
	"github.com/penguin-statistics/widget-backend/utils"
	"github.com/sirupsen/logrus"
)

// Controller consists instance for a type of data
type Controller struct {
	caches []*utils.Cache
	logger *logrus.Entry
}

// Stage specifies data structure for the stage data type
type Stage struct {
	ZoneID       string            `json:"zoneId,omitempty"`
	StageID      string            `json:"stageId,omitempty"`
	CodeI18N     map[string]string `json:"code_i18n,omitempty"`
	ApCost       int               `json:"apCost,omitempty"`
	Existence    models.Existence  `json:"existence"`
	MinClearTime int               `json:"minClearTime,omitempty"`
}

func (s Stage) MarshalJSON() ([]byte, error) {
	var tmp struct {
		ZoneID       string            `json:"zoneId,omitempty"`
		StageID      string            `json:"stageId,omitempty"`
		CodeI18N     map[string]string `json:"code_i18n,omitempty"`
		ApCost       int               `json:"apCost,omitempty"`
		MinClearTime int               `json:"minClearTime,omitempty"`
	}
	tmp.ZoneID = s.ZoneID
	tmp.StageID = s.StageID
	tmp.CodeI18N = s.CodeI18N
	tmp.ApCost = s.ApCost
	tmp.MinClearTime = s.MinClearTime
	return json.Marshal(&tmp)
}
