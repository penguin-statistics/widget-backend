package errors

import (
	"net/http"
)

const (
	// BlameUser describes the error was caused by user and the status code shall blame the user
	BlameUser   = http.StatusBadRequest
	// BlameServer describes the error was caused by server internally and therefore shall blame the server
	BlameServer = http.StatusInternalServerError
)

// Error describes an error with several metadata to describe such error
type Error struct {
	Type    string      `json:"type"`
	Details interface{} `json:"details"`
	Blame   int         `json:"-"`
}

// WrappedError wraps Error for the client to be able to easily distinguish successful responses and error responses
type WrappedError struct {
	Error *Error `json:"error"`
}

// New creates a new Error instance
func New(typ string, details interface{}, blame int) *Error {
	return &Error{
		Type:    typ,
		Details: details,
		Blame:   blame,
	}
}

// Wrapped returns a WrappedError consists of current Error
func (e *Error) Wrapped() *WrappedError {
	return &WrappedError{Error: e}
}
