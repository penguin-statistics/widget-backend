package utils

import (
	"net/http"
	"net/url"
)

const (
	// BaseEndpoint is the endpoint for all of the data
	BaseEndpoint = "https://penguin-stats.io/PenguinStats/api/v2/"
)

// NewPenguinRequest creates a http.Request that parses resource as relative path against BaseEndpoint to produce path, and appends server to the URL that is being generated
func NewPenguinRequest(resource string, params *url.Values) (*http.Request, error) {
	rel, err := url.ParseRequestURI(BaseEndpoint)
	if err != nil {
		return nil, err
	}
	rel, err = rel.Parse(resource)
	if err != nil {
		return nil, err
	}
	rel.RawQuery = params.Encode()

	r, err := http.NewRequest("GET", rel.String(), nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "PenguinStatsWidgetBackend/1.0")
	return r, nil
}
