package zone

import (
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
	ZoneNameI18N map[string]string `json:"zoneName_i18n,omitempty"`
}
