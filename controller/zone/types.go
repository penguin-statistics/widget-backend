package zone

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

// Zone specifies data structure for the zone data type
type Zone struct {
	ZoneID       string            `json:"zoneId,omitempty"`
	Type         string            `json:"type,omitempty"`
	Existence    models.Existence  `json:"existence"`
	ZoneNameI18N map[string]string `json:"zoneName_i18n,omitempty"`
}

func (z Zone) MarshalJSON() ([]byte, error) {
	var tmp struct {
		ZoneID       string            `json:"zoneId,omitempty"`
		Type         string            `json:"type,omitempty"`
		ZoneNameI18N map[string]string `json:"zoneName_i18n,omitempty"`
	}
	tmp.ZoneID = z.ZoneID
	tmp.Type = z.Type
	tmp.ZoneNameI18N = z.ZoneNameI18N
	return json.Marshal(&tmp)
}
