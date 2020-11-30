package response

import (
	"github.com/penguin-statistics/widget-backend/controller/item"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/stage"
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/controller/zone"
)

// MatrixResponse consists additional data from the matrix result itself
type MatrixResponse struct {
	Query *matrix.Query `json:"query,omitempty"`

	// CacheStatus represents statuses of underlying controllers and caches of current response
	CacheStatus map[string]*status.Status `json:"cache,omitempty"`

	Items  []*item.Item     `json:"items,omitempty"`
	Matrix []*matrix.Matrix `json:"matrix,omitempty"`
	Stages []*stage.Stage   `json:"stages,omitempty"`
	Zones  []*zone.Zone     `json:"zones,omitempty"`
}

// NewResponse creates a new MatrixResponse
func NewResponse() *MatrixResponse {
	return &MatrixResponse{
		CacheStatus: map[string]*status.Status{},

		Items:  []*item.Item{},
		Matrix: []*matrix.Matrix{},
		Stages: []*stage.Stage{},
		Zones:  []*zone.Zone{},
	}
}
