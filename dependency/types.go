package dependency

import (
	"github.com/penguin-statistics/partial-matrix/controller/item"
	"github.com/penguin-statistics/partial-matrix/controller/matrix"
	"github.com/penguin-statistics/partial-matrix/controller/stage"
	"github.com/penguin-statistics/partial-matrix/controller/zone"
)

type MatrixResponse struct {
	Problems map[string]string `json:"problems,omitempty"`

	Items  []*item.Item     `json:"items,omitempty"`
	Matrix []*matrix.Matrix `json:"matrix,omitempty"`
	Stages []*stage.Stage   `json:"stages,omitempty"`
	Zones  []*zone.Zone     `json:"zones,omitempty"`
}

func NewResponse() *MatrixResponse {
	return &MatrixResponse{
		Problems: map[string]string{},

		Items:  []*item.Item{},
		Matrix: []*matrix.Matrix{},
		Stages: []*stage.Stage{},
		Zones:  []*zone.Zone{},
	}
}
