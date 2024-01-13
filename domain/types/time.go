package types

import (
	"fmt"
	"strings"
	"time"
)

const ElionaTimeFormat = "2006-01-02T15:04:05.999Z07:00"

// ElionaTimestamp is a wrapper for time.Time ensuring time format used in Eliona
type ElionaTimestamp struct {
	time.Time
}

func NewElionaTimestamp() ElionaTimestamp {
	return ElionaTimestamp{
		Time: time.Now(),
	}
}

func (et *ElionaTimestamp) MarshalJSON() ([]byte, error) {
	if et.IsZero() {
		return nil, fmt.Errorf("nil time") // TODO: real error
	}

	return []byte(et.Format(ElionaTimeFormat)), nil
}

func (et *ElionaTimestamp) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		et.Time = time.Time{}
		return
	}

	et.Time, err = time.Parse(ElionaTimeFormat, s)

	return err
}
