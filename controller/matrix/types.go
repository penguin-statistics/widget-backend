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
	Start    *int64   `json:"start,omitempty"`
	End      *int64   `json:"end,omitempty"`
}

// Query describes a query on Matrix
type Query struct {
	// StageID is the query stageId of the current response; empty represents stageId hasn't been used as one of the constraints
	StageID string `json:"stageId,omitempty"`

	// ItemID is the query itemId of the current response; empty represents itemId hasn't been used as one of the constraints
	ItemID string `json:"itemId,omitempty"`

	// Server is the server to query matrix of
	Server string `json:"server,omitempty"`
}
