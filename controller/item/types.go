package item

import (
	"github.com/penguin-statistics/partial-matrix/utils"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	cache  *utils.Cache
	logger *logrus.Entry
}

type Item struct {
	ItemID      string            `json:"itemId,omitempty"`
	NameI18N    map[string]string `json:"name_i18n,omitempty"`
	SpriteCoord []int             `json:"spriteCoord,omitempty"`
}
