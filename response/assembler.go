package response

import (
	"bytes"
	"fmt"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/meta"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/utils"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
)

// Assembler is to marshal records
type Assembler struct {
	collection *meta.Collection
	tmpl       *template.Template
}

// New creates a new Assembler that marshals records
func New(collection *meta.Collection, resourceLocation string) *Assembler {
	tmpl, err := ioutil.ReadFile(path.Join(resourceLocation, "/index.html"))
	if err != nil {
		panic(err)
	}

	return &Assembler{
		collection: collection,
		tmpl:       template.Must(template.New("widget").Parse(string(tmpl))),
	}
}

// Marshal marshals records with their dependencies and gives MatrixResponse that contains rich metadata for current state of controllers
func (m *Assembler) Marshal(records []*matrix.Matrix, server string, query *MatrixQuery) (*MatrixResponse, error) {
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
	response.Query = query
	response.CacheStatus = m.collection.Statuses(server)

	return response, nil
}

func (m *Assembler) Response(response *MatrixResponse) []byte {
	buf := bytes.Buffer{}
	err := m.tmpl.Execute(&buf, struct {
		PenguinWidgetData *MatrixResponse
	}{
		PenguinWidgetData: response,
	})
	if err != nil {
		return []byte(fmt.Sprintf("data injection failed with error %v", err))
	}
	return buf.Bytes()
}

func (m *Assembler) Error(error *errors.Error) (int, []byte) {
	buf := bytes.Buffer{}
	err := m.tmpl.Execute(&buf, struct {
		PenguinWidgetData *errors.WrappedError
	}{
		PenguinWidgetData: error.Wrapped(),
	})
	if err != nil {
		return http.StatusInternalServerError, []byte(fmt.Sprintf("fatal: data injection failed with error %v", err))
	}
	return error.Blame, buf.Bytes()
}
