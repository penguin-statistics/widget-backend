package errors

import (
	"encoding/json"
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

func (e *Error) Error() string {
	b, err := json.Marshal(&e)
	if err != nil {
		return "(marshalling failed) " + e.Type
	}
	return string(b)
}
