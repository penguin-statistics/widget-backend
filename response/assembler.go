package response

import (
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/meta"
	"github.com/penguin-statistics/widget-backend/utils"
)

// Assembler is to marshal records
type Assembler struct {
	collection *meta.Collection
}

// New creates a new Assembler that marshals records
func New(collection *meta.Collection) *Assembler {
	return &Assembler{collection: collection}
}

// Marshal marshals records with their dependencies and gives MatrixResponse that contains rich metadata for current state of controllers
func (m *Assembler) Marshal(records []*matrix.Matrix, server string, typ string) (*MatrixResponse, error) {
	response := NewResponse()

	itemDeps := utils.NewUniString()
	stageDeps := utils.NewUniString()
	zoneDeps := utils.NewUniString()

	// dependency collection and injection
	{
		for _, record := range records {
			itemDeps.Add(record.ItemID)
			stageDeps.Add(record.StageID)
		}
		for _, stageID := range stageDeps.Slice() {
			stage, err := m.collection.Stage.Stage(stageID)
			if err != nil {
				return nil, err
			}
			response.Stages = append(response.Stages, stage)
			zoneDeps.Add(stage.ZoneID)
		}
		for _, itemID := range itemDeps.Slice() {
			item, err := m.collection.Item.Item(itemID)
			if err != nil {
				return nil, err
			}
			response.Items = append(response.Items, item)
		}
		for _, zoneID := range zoneDeps.Slice() {
			zone, err := m.collection.Zone.Zone(zoneID)
			if err != nil {
				return nil, err
			}
			response.Zones = append(response.Zones, zone)
		}
	}

	// populate remaining fields
	response.Matrix = records
	response.Type = typ
	response.CacheStatus = m.collection.Statuses(server)

	return response, nil
}
