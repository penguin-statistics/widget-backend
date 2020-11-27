package dependency

import (
	"github.com/penguin-statistics/partial-matrix/controller"
	"github.com/penguin-statistics/partial-matrix/controller/matrix"
	"github.com/penguin-statistics/partial-matrix/utils"
)

type Manager struct {
	collection *controller.Collection
}

func New(collection *controller.Collection) *Manager {
	return &Manager{collection: collection}
}

func (m *Manager) Populate(records []*matrix.Matrix) (*MatrixResponse, error) {
	response := NewResponse()

	itemDeps := utils.NewUniString()
	stageDeps := utils.NewUniString()
	zoneDeps := utils.NewUniString()

	for _, record := range records {
		itemDeps.Add(record.ItemID)
		stageDeps.Add(record.StageID)
	}
	for _, stageId := range stageDeps.Slice() {
		stage, err := m.collection.Stage.Stage(stageId)
		if err != nil {
			return nil, err
		}
		response.Stages = append(response.Stages, stage)
		zoneDeps.Add(stage.ZoneID)
	}
	for _, itemId := range itemDeps.Slice() {
		item, err := m.collection.Item.Item(itemId)
		if err != nil {
			return nil, err
		}
		response.Items = append(response.Items, item)
	}
	for _, zoneId := range zoneDeps.Slice() {
		zone, err := m.collection.Zone.Zone(zoneId)
		if err != nil {
			return nil, err
		}
		response.Zones = append(response.Zones, zone)
	}

	response.Matrix = records

	return response, nil
}
