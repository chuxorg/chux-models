package models

import (
	"encoding/json"
	"time"
)

// CustomTime is a struct used to hold time types in a struct
// so that the struct can be marshalled, unmarshalled
type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	var dateString string
	err := json.Unmarshal(b, &dateString)
	if err != nil {
		return err
	}

	dateFormats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	var parsedTime time.Time
	for _, format := range dateFormats {
		parsedTime, err = time.Parse(format, dateString)
		if err == nil {
			ct.Time = parsedTime
			return nil
		}
	}

	return err
}

func (ct *CustomTime) Now(){
	
	ct.Time = time.Now()
}
