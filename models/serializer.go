package models

import "encoding/json"

// Marshalls a struct into a slice of bytes
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal takes JSON data as bytes and an interface value to unmarshal the data into,
// and returns the un-marshalled interface value.
func Unmarshal(data []byte, v interface{}) (interface{}, error) {
	err := json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
