package kafka

import (
	"encoding/json"

	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	dlqTopic = "eliona.dlq.log.v1"
)

var dlqKey = []byte("jobsync")

type dlqReport struct {
	Key   string
	Value string
	Error string
}

func (c *Client) ReportDLQ(key, value, err string) {
	report := dlqReport{
		Key:   key,
		Value: value,
		Error: err,
	}

	b, _ := json.Marshal(report)

	c.client.Produce(c.ctx, &kgo.Record{
		Topic: dlqTopic,
		Key:   dlqKey,
		Value: b,
	}, c.produceCallback)
}
