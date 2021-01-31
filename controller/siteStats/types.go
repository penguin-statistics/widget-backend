package siteStats

import (
	"github.com/penguin-statistics/widget-backend/controller/item"
	"github.com/penguin-statistics/widget-backend/controller/stage"
	"github.com/penguin-statistics/widget-backend/utils"
	"github.com/sirupsen/logrus"
)

// Controller consists instance for a type of data. SiteStat controller consists a slice of utils.Cache which contains data from multiple server
type Controller struct {
	caches []*utils.Cache
	logger *logrus.Entry
}

// SiteStat specifies data structure for the site stats data type
type SiteStat struct {
	Stage  *stage.Stage `json:"stage,omitempty"`
	Item   *item.Item `json:"item,omitempty"`
	Quantity int    `json:"quantity"` // cannot omitempty as it is zeroable
	Times    int    `json:"times"`    // cannot omitempty as it is zeroable
	RecentTimes    int    `json:"recentTimes"`    // cannot omitempty as it is zeroable
}

type RemoteResponse struct {
	TotalStageTimesRecent24H  []StageTime `json:"totalStageTimes_24h,omitempty"`
}

type StageTime struct {
	StageID     string `json:"stageId,omitempty"`
	RecentTimes int    `json:"times"`
}

// Query describes a query on SiteStat
type Query struct {
	// Server is the server to query matrix of
	Server string `json:"server,omitempty"`
}
