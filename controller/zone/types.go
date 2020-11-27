package zone

import (
	"github.com/penguin-statistics/partial-matrix/utils"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	cache  *utils.Cache
	logger *logrus.Entry
}

type Zone struct {
	ZoneID       string            `json:"zoneId,omitempty"`
	Type         string            `json:"type,omitempty"`
	ZoneNameI18N map[string]string `json:"zoneName_i18n,omitempty"`
}
