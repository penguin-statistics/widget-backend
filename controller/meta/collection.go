package meta

import (
	"github.com/penguin-statistics/widget-backend/controller/item"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/stage"
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/controller/zone"
)

// Collection consists of all local controllers that manages records
type Collection struct {
	Item   *item.Controller
	Matrix *matrix.Controller
	Stage  *stage.Controller
	Zone   *zone.Controller
}

// NewCollection creates all underlying controllers and returns Collection
func NewCollection() *Collection {
	return &Collection{
		Item:   item.New(),
		Matrix: matrix.New(),
		Stage:  stage.New(),
		Zone:   zone.New(),
	}
}

// Statuses collects cache status from its underlying utils.Cache instance and return a map which its key is the controller name, and value is the status.Status
func (c *Collection) Statuses(server string) map[string]*status.Status {
	return map[string]*status.Status{
		"item":   c.Item.Status(server),
		"matrix": c.Matrix.Status(server),
		"stage":  c.Stage.Status(server),
		"zone":   c.Zone.Status(server),
	}
}
