package kafka

import (
	"context"
	"errors"

	"github.com/eliona-smart-building-assistant/go-utils/log"
	"github.com/insolare/elionacore/domain/types"
)

// HandleFunc is a handler called for each message received over Kafka.
// If returned error is nil - offset is marked as committed.
type HandleFunc func(m types.Message) error

func (c *Client) consumeSimple(handler HandleFunc) {
	c.logger.Info(c.facility, "Simple consumer")
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

			err := handler(m)
			if err != nil {
				log.Error(c.facility, "Error in handler function: %v", err)
			}

			c.client.MarkCommitRecords(rec)
		}
	}
}
