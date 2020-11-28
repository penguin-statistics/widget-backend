package status

import "time"

// Status describes status of current cache
type Status struct {
	UpdatedAt *time.Time `json:"updated"`
	FailCount int        `json:"fails,omitempty"`
	// Length is not really useful when returning a response, therefore we omit it
	Length int `json:"-"`
}
