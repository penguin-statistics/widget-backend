package response

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/meta"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/utils"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"time"
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
func (m *Assembler) Marshal(records []*matrix.Matrix, query *matrix.Query) *MatrixResponse {
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
			stage := m.collection.Stage.Stage(query.Server, stageID)
			response.Stages = append(response.Stages, stage)
			zoneDeps.Add(stage.ZoneID)
		}
		for _, itemID := range itemDeps.Slice() {
			item := m.collection.Item.Item(itemID)
			response.Items = append(response.Items, item)
		}
		for _, zoneID := range zoneDeps.Slice() {
			zone := m.collection.Zone.Zone(zoneID)
			response.Zones = append(response.Zones, zone)
		}
	}

	// populate remaining fields
	response.Matrix = records
	response.Query = query
	response.CacheStatus = m.collection.Statuses(query.Server)

	return response
}

// inject injects Last-Modified headers along with other metadata that is used
func inject(c echo.Context, response *MatrixResponse) *errors.Error {
	// == Last-Modified
	lastModified, err := time.Parse("2006-01-02", "1970-01-01")
	if err != nil {
		return errors.New("InternalError", "failed to calculate Last-Modified header: initialize failed", errors.BlameServer)
	}
	initialTime := lastModified
	for _, status := range response.CacheStatus {
		// if current is *later* (#After) then use this instead of old one
		if status.UpdatedAt.After(lastModified) {
			lastModified = *status.UpdatedAt
		}
	}
	if lastModified == initialTime {
		return errors.New("InternalError", "failed to calculate Last-Modified header: malformed cache data", errors.BlameServer)
	}

	// last modified according of the current cache status
	c.Response().Header().Add(echo.HeaderLastModified, lastModified.UTC().Format(time.RFC1123))
	// indicate Vary to improve cache behavior
	c.Response().Header().Add("Vary", "CF-IPCountry")

	// RequestMetadata

	metadata := c.Get("meta").(*RequestMetadata)
	response.Request = metadata

	return nil
}

func (m *Assembler) HTMLResponse(c echo.Context, response *MatrixResponse) error {
	// inject Last-Modified headers and other metadata
	injectErr := inject(c, response)
	if injectErr != nil {
		return c.HTMLBlob(m.HTMLError(injectErr))
	}

	buf := bytes.Buffer{}
	err := m.tmpl.Execute(&buf, struct {
		PenguinWidgetData *MatrixResponse
	}{
		PenguinWidgetData: response,
	})
	if err != nil {
		return c.HTMLBlob(m.HTMLError(errors.New("CantMarshal", "data injection failed", errors.BlameServer)))
	}
	body := buf.Bytes()

	return c.HTMLBlob(http.StatusOK, body)
}

func (m *Assembler) JSONResponse(c echo.Context, response *MatrixResponse) error {
	// inject Last-Modified headers
	injectErr := inject(c, response)
	if injectErr != nil {
		return c.JSON(injectErr.Blame, injectErr.Wrapped())
	}

	return c.JSON(http.StatusOK, response)
}

func (m *Assembler) HTMLError(error *errors.Error) (int, []byte) {
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

func (m *Assembler) JSONError(error *errors.Error) (int, interface{}) {
	return error.Blame, error.Wrapped()
}

