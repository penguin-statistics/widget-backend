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
func NewPenguinRequest(resource string, server string) (*http.Request, error) {
	rel, err := url.ParseRequestURI(BaseEndpoint)
	if err != nil {
		return nil, err
	}
	rel, err = rel.Parse(resource)
	if err != nil {
		return nil, err
	}
	if server != "" {
		q := rel.Query()
		q.Set("server", server)
		rel.RawQuery = q.Encode()
	}

	return http.NewRequest("GET", rel.String(), nil)
}
