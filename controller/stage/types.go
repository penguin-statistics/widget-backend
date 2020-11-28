package stage

import (
	"github.com/penguin-statistics/widget-backend/utils"
	"github.com/sirupsen/logrus"
)

// Controller consists instance for a type of data
type Controller struct {
	cache  *utils.Cache
	logger *logrus.Entry
}

// Stage specifies data structure for the stage data type
type Stage struct {
	ZoneID       string            `json:"zoneId,omitempty"`
	StageID      string            `json:"stageId,omitempty"`
	CodeI18N     map[string]string `json:"code_i18n,omitempty"`
	ApCost       int               `json:"apCost,omitempty"`
	MinClearTime int               `json:"minClearTime,omitempty"`
}
