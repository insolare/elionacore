package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// HandleFunc is called on each event on listenned channels. 1st argument is channel name and second - the event itself
type HandleFunc func(string, string)

func (c *Client) EventListener(fn HandleFunc, channels ...string) {
	c.wg.Add(1)
	go c.listener(fn, channels...)
}

func (c *Client) listener(fn HandleFunc, channels ...string) {
	defer c.wg.Done()

	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	conn, err := c.pool.Acquire(ctx)
	cancel()

	if err != nil {
		c.errorCh <- fmt.Errorf("error acquiring connection: %v", err)
		return
	}

	defer conn.Release()

	for _, channel := range channels {
		ctx, cancel := context.WithTimeout(c.ctx, 1*time.Second)
		query := "LISTEN " + channel
		_, err := conn.Exec(ctx, query)
		cancel()

		if err != nil {
			c.errorCh <- fmt.Errorf("error subscribing to events on %s: %v", channel, err)
			return
		}
	}

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.stop:
			return
		default:
		}

		ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
		n, err := conn.Conn().WaitForNotification(ctx)
		cancel()
		if err != nil {
			if err != context.Canceled && !pgconn.Timeout(err) {
				c.errorCh <- fmt.Errorf("critical error while listening for events: %v", err)
				return
			}
			continue
		}

		if fn != nil {
			fn(n.Channel, n.Payload)
		}
	}
}
