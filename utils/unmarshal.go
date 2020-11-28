package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// UnmarshalFromReader unmarshals all data from reader and populates it to v
func UnmarshalFromReader(reader io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
