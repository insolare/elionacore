package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type sampleStruct struct {
	Timestamp ElionaTimestamp `json:"ts"`
}

func TestTimestampUnmarhsal(t *testing.T) {
	assert := assert.New(t)

	body := `{"ts": "2022-12-31T23:01:02.123+03:00"}`

	valid, _ := time.Parse(ElionaTimeFormat, "2022-12-31T23:01:02.123+03:00")
	check := sampleStruct{}
	err := json.Unmarshal([]byte(body), &check)

	assert.Nil(err)
	assert.True(valid.Equal(check.Timestamp.Time))
}

func TestTimestampMarshal(t *testing.T) {
	assert := assert.New(t)
	ts := NewElionaTimestamp()

	s := sampleStruct{
		Timestamp: ts,
	}

	b, err := json.Marshal(s)
	assert.Nil(err)

	assert.Equal("", string(b))
}
