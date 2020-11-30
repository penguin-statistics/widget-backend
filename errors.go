package main

import (
	"github.com/penguin-statistics/widget-backend/errors"
)

var (
	ErrInvalidServer = errors.New("InvalidServer", "malformed parameter `server` provided", errors.BlameUser)
	ErrCantMarshal   = errors.New("CantMarshal", "failed to populate matrix with metadata provided", errors.BlameServer)
)
