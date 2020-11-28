package item

import (
	"github.com/penguin-statistics/widget-backend/utils"
	"github.com/sirupsen/logrus"
)

// Controller consists instance for a type of data
type Controller struct {
	cache  *utils.Cache
	logger *logrus.Entry
}

// Item specifies data structure for the item data type
type Item struct {
	ItemID      string            `json:"itemId,omitempty"`
	NameI18N    map[string]string `json:"name_i18n,omitempty"`
	SpriteCoord []int             `json:"spriteCoord,omitempty"`
}
