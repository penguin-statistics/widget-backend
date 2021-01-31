package item

import (
	"encoding/json"
	"github.com/penguin-statistics/widget-backend/models"
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
	Existence   models.Existence  `json:"existence"`
	SpriteCoord []int             `json:"spriteCoord,omitempty"`
}

func (i Item) MarshalJSON() ([]byte, error) {
	var tmp struct {
		ItemID      string            `json:"itemId,omitempty"`
		NameI18N    map[string]string `json:"name_i18n,omitempty"`
		SpriteCoord []int             `json:"spriteCoord,omitempty"`
	}
	tmp.ItemID = i.ItemID
	tmp.NameI18N = i.NameI18N
	tmp.SpriteCoord = i.SpriteCoord
	return json.Marshal(&tmp)
}
