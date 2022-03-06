package response

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"sort"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/meta"
	"github.com/penguin-statistics/widget-backend/controller/siteStats"
	"github.com/penguin-statistics/widget-backend/errors"
	"github.com/penguin-statistics/widget-backend/utils"
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

// MarshalMatrix marshals records with their dependencies and gives MatrixResponse that contains rich metadata for current state of controllers
func (m *Assembler) MarshalMatrix(records []*matrix.Matrix, query *matrix.Query) *MatrixResponse {
	response := NewMatrixResponse()

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

// MarshalSiteStats marshals records with their dependencies and gives MatrixResponse that contains rich metadata for current state of controllers
func (m *Assembler) MarshalSiteStats(records []siteStats.StageTime, query *siteStats.Query) *SiteStatsResponse {
	stats := []*siteStats.SiteStat{}

OUTER:
	for _, record := range records {
		matrices, err := m.collection.Matrix.Query(&matrix.Query{Server: query.Server, StageID: record.StageID})
		if err != nil {
			continue
		}
		sort.SliceStable(matrices, func(i, j int) bool {
			return (float64(matrices[i].Quantity) / float64(matrices[i].Times)) > (float64(matrices[j].Quantity) / float64(matrices[j].Times))
		})

		now := time.Now()
		for _, mat := range matrices {
			if mat.Start != nil {
				start := time.Unix(*mat.Start/1000, 0)
				// not yet arrived
				if start.After(now) {
					continue
				}
			}

			if mat.End != nil {
				end := time.Unix(*mat.End/1000, 0)
				// already passed
				if end.Before(now) {
					continue
				}
			}

			foundStage := m.collection.Stage.Stage(query.Server, mat.StageID)
			// stage & item shall always qualify: right after the activity ends, server shall still allow
			// some of the records to be able to show on the widget side

			foundItem := m.collection.Item.Item(mat.ItemID)

			stat := &siteStats.SiteStat{
				Stage:       foundStage,
				Item:        foundItem,
				Quantity:    mat.Quantity,
				Times:       mat.Times,
				RecentTimes: record.RecentTimes,
			}

			stats = append(stats, stat)
			continue OUTER
		}
	}

	response := &SiteStatsResponse{}
	// populate remaining fields
	response.Stats = stats
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

// HTMLResponse renders the response into HTML format and writes it into c
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

// JSONMatrixResponse renders the response into JSON format and writes it into c
func (m *Assembler) JSONMatrixResponse(c echo.Context, response *MatrixResponse) error {
	// inject Last-Modified headers
	injectErr := inject(c, response)
	if injectErr != nil {
		return c.JSON(injectErr.Blame, injectErr.Wrapped())
	}

	return c.JSON(http.StatusOK, response)
}

// HTMLError renders the error into HTML format and writes it into c
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

// JSONError renders the error into JSON format and writes it into c
func (m *Assembler) JSONError(error *errors.Error) (int, interface{}) {
	return error.Blame, error.Wrapped()
}
