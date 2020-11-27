package utils

import (
	"net/http"
	"net/url"
)

const (
	BaseEndpoint = "https://penguin-stats.io/PenguinStats/api/v2/"
)

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
