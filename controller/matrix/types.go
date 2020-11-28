package matrix

import (
	"github.com/penguin-statistics/widget-backend/utils"
	"github.com/sirupsen/logrus"
)

// Controller consists instance for a type of data. Matrix controller consists a slice of utils.Cache which contains data from multiple server
type Controller struct {
	caches []*utils.Cache
	logger *logrus.Entry
}

// Matrix specifies data structure for the matrix data type
type Matrix struct {
	StageID  string `json:"stageId,omitempty"`
	ItemID   string `json:"itemId,omitempty"`
	Quantity int    `json:"quantity"` // cannot omitempty as it is zeroable
	Times    int    `json:"times"`    // cannot omitempty as it is zeroable
	Start    *int   `json:"start,omitempty"`
	End      *int   `json:"end,omitempty"`
}
