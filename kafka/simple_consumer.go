package kafka

import (
	"context"
	"errors"

	"github.com/insolare/elionacore/domain/types"
)

// HandleFunc is a handler called for each message received over Kafka.
// If returned error is nil - offset is marked as committed.
type HandleFunc func(*Client, types.Message)

func (c *Client) consumeSimple(handler HandleFunc) {
	c.logger.Info(c.facility, "Simple consumer started")
	defer c.logger.Info(c.facility, "Simple consumer stopped")

	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.stop:
			return
		default:
		}

		fetches := c.client.PollRecords(c.ctx, 1000)
		if fetches.IsClientClosed() || errors.Is(fetches.Err0(), context.Canceled) {
			return
		}

		//  TODO: Check errors for each!

		iter := fetches.RecordIter()
		for !iter.Done() {
			rec := iter.Next()

			m := types.Message{
				Topic: rec.Topic,
				Key:   rec.Key,
				Value: rec.Value,
			}

			handler(c, m)

			c.client.MarkCommitRecords(rec)
		}
	}
}
