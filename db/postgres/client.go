package postgres

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	ctx       context.Context
	cancel    context.CancelFunc
	pool      *pgxpool.Pool
	errorCh   chan error
	wg        sync.WaitGroup
	stop      chan struct{}
	opTimeout time.Duration
}

type ClientConfig struct {
	DSN              string
	AppName          string
	OperationTimeout time.Duration
}

func NewClient(ctx context.Context, cfg ClientConfig, errorCh chan error) (*Client, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.Config.RuntimeParams["application_name"] = cfg.AppName

	connCtx, connCancel := context.WithTimeout(ctx, cfg.OperationTimeout)
	defer connCancel()

	pool, err := pgxpool.NewWithConfig(connCtx, poolConfig)
	if err != nil {
		return nil, err
	}

	clientCtx, clientCancel := context.WithCancel(ctx)
	db := &Client{
		ctx:       clientCtx,
		cancel:    clientCancel,
		pool:      pool,
		errorCh:   errorCh,
		wg:        sync.WaitGroup{},
		stop:      make(chan struct{}),
		opTimeout: cfg.OperationTimeout,
	}

	return db, nil
}

func (c *Client) Close() {
	close(c.stop)
	c.pool.Close()
	c.wg.Wait()
}

// Multiple structs, value maps by position
func QueryStructSliceByPos[T any](c *Client, query string, params ...any) ([]T, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.opTimeout)
	defer cancel()

	rows, err := c.pool.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return pgx.CollectRows[T](rows, pgx.RowToStructByPos[T])
}

// Single scalar result
func QueryOne[T any](c *Client, query string, params ...any) (T, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.opTimeout)
	defer cancel()

	var result T
	err := c.pool.QueryRow(ctx, query, params...).Scan(&result)

	return result, err
}

// Multiple scalar results
func QueryMany[T any](c *Client, query string, params ...any) ([]T, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.opTimeout)
	defer cancel()

	rows, err := c.pool.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]T, 0, 1)
	for rows.Next() {
		var nxt T
		err = rows.Scan(&nxt)
		if err != nil {
			continue
		}

		results = append(results, nxt)
	}

	return results, nil
}
