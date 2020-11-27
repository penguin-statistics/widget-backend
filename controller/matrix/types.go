package matrix

import (
	"github.com/penguin-statistics/partial-matrix/utils"
	"github.com/sirupsen/logrus"
)

type Controller struct {
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
