package response

import (
	"github.com/penguin-statistics/widget-backend/controller/item"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/siteStats"
	"github.com/penguin-statistics/widget-backend/controller/stage"
	"github.com/penguin-statistics/widget-backend/controller/status"
	"github.com/penguin-statistics/widget-backend/controller/zone"
)

// RequestMetadata describes metadata related to the subsequent request
type RequestMetadata struct {
	// Mirror is the preferred mirror to select from; oftenly chose with reference of `CF-IPCountry` header
	Mirror string `json:"mirror,omitempty"`
}

// SiteStatsResponse consists additional data from the site stats result itself
type SiteStatsResponse struct {
	Query *siteStats.Query `json:"query,omitempty"`

	// CacheStatus represents statuses of underlying controllers and caches of current response
	CacheStatus map[string]*status.Status `json:"cache,omitempty"`

	Stats  []*siteStats.SiteStat     `json:"stats,omitempty"`
}

// MatrixResponse consists additional data from the matrix result itself
type MatrixResponse struct {
	Request *RequestMetadata `json:"request,omitempty"`
	Query *matrix.Query `json:"query,omitempty"`

	// CacheStatus represents statuses of underlying controllers and caches of current response
	CacheStatus map[string]*status.Status `json:"cache,omitempty"`

	Items  []*item.Item     `json:"items,omitempty"`
	Matrix []*matrix.Matrix `json:"matrix,omitempty"`
	Stages []*stage.Stage   `json:"stages,omitempty"`
	Zones  []*zone.Zone     `json:"zones,omitempty"`
}

// NewMatrixResponse creates a new MatrixResponse
func NewMatrixResponse() *MatrixResponse {
	return &MatrixResponse{
		CacheStatus: map[string]*status.Status{},

		Items:  []*item.Item{},
		Matrix: []*matrix.Matrix{},
		Stages: []*stage.Stage{},
		Zones:  []*zone.Zone{},
	}
}
