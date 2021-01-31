package models

import "time"

type Existence map[string]ServerExistence
type ServerExistence struct {
	Exist bool `json:"exist"`
	OpenTime *int64 `json:"openTime"`
	CloseTime *int64 `json:"closeTime"`
}

// IsExist checks an Existence instance to see if it is existed in server
func (e Existence) IsExist(server string) bool {
	existence, ok := e[server]
	if !ok {
		// if no such server we assume it does not exist
		return false
	}
	if !existence.Exist {
		// explicitly not exist
		return false
	}

	now := time.Now()

	/*
	   Invalid: open.After(now) == true      Valid: !open.After(now)
	     |                                    |  && !closed.Before(now)
	     |                                    |
	     v                                    v
	   |||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||
	        ^                                                       ^
	       OPEN                                                   CLOSE
	 */

	// only consider OpenTime if it is provided
	if existence.OpenTime != nil {
		open := time.Unix(*existence.OpenTime / 1000, 0)
		// open is after now: not yet opened
		if open.After(now) {
			return false
		}
	}

	// only consider CloseTime if it is provided
	if existence.CloseTime != nil {
		closed := time.Unix(*existence.CloseTime / 1000, 0)
		// closed is before now: already closed
		if closed.Before(now) {
			return false
		}
	}

	// all conditions met
	return true
}