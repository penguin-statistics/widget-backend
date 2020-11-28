package errors

import (
	"net/http"
)

const (
	BlameUser   = http.StatusBadRequest
	BlameServer = http.StatusInternalServerError
)

type Error struct {
	Type    string      `json:"type"`
	Details interface{} `json:"details"`
	Blame   int         `json:"-"`
}

type WrappedError struct {
	Error *Error `json:"error"`
}

func New(typ string, details interface{}, blame int) *Error {
	return &Error{
		Type:    typ,
		Details: details,
		Blame:   blame,
	}
}

func (e *Error) Wrapped() *WrappedError {
	return &WrappedError{Error: e}
}
