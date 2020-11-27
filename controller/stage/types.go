package stage

import (
	"github.com/penguin-statistics/partial-matrix/utils"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	cache  *utils.Cache
	logger *logrus.Entry
}

type Stage struct {
	ZoneID       string            `json:"zoneId,omitempty"`
	StageID      string            `json:"stageId,omitempty"`
	CodeI18N     map[string]string `json:"code_i18n,omitempty"`
	ApCost       int               `json:"apCost,omitempty"`
	MinClearTime int               `json:"minClearTime,omitempty"`
}
