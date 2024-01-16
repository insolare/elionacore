package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/insolare/elionacore/domain/types"
	"github.com/insolare/elionacore/tools"
	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	facility = "kafka"
)

type Client struct {
	ctx      context.Context
	cancel   context.CancelFunc
	client   *kgo.Client
	wg       sync.WaitGroup
	closing  bool
	stop     chan struct{}
	logger   types.Logger
	facility string
}

type ClientConfig struct {
	Seeds       []string
	Group       string
	ClientID    string
	Logger      types.Logger
	HandlerFn   HandleFunc
	HandleAsync bool
}

var globalClient *Client

// TODO: Make initialization of _default_ client
func Init() {}

func NewClient(ctx context.Context, cfg ClientConfig) (*Client, error) {
	opts := make([]kgo.Opt, 0)
	opts = append(opts, kgo.BlockRebalanceOnPoll(), kgo.AutoCommitMarks())

	var logger types.Logger
	if cfg.Logger == nil {
		logger = types.NoopLogger{}
	} else {
		logger = cfg.Logger
	}

	if len(cfg.Seeds) < 1 {
		return nil, ErrNoSeeds
	}

	opts = append(opts, kgo.SeedBrokers(cfg.Seeds...))

	if len(cfg.ClientID) < 1 {
		cfg.ClientID = tools.GetRandomName(0)
		logger.Warning(facility, "No ClientID provided in config. Will use '%s' as ClientID and logging facility")
	}

	opts = append(opts, kgo.ClientID(cfg.ClientID))

	if len(cfg.Group) > 0 {
		opts = append(opts, kgo.ConsumerGroup(cfg.Group))
	}

	franz, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	err = franz.Ping(pingCtx)
	if err != nil {
		return nil, err
	}

	clientCtx, cancel := context.WithCancel(ctx)
	client := &Client{
		ctx:      clientCtx,
		cancel:   cancel,
		client:   franz,
		wg:       sync.WaitGroup{},
		stop:     make(chan struct{}),
		logger:   logger,
		facility: facility + ":" + cfg.ClientID,
	}

	// TODO: Select actual consumer type here
	client.wg.Add(1)
	go client.consumeSimple(messageHandler(cfg.HandlerFn, cfg.HandleAsync))

	return client, nil
}

func messageHandler(fn HandleFunc, async bool) HandleFunc {
	if async {
		return func(m types.Message) {
			go fn(m)
		}
	}

	return func(m types.Message) {
		fn(m)
	}
}

func Subscribe(topics ...string) {
	if globalClient != nil {
		globalClient.Subscribe(topics...)
	}
}

func Unsubscribe(topics ...string) {
	if globalClient != nil {
		globalClient.Unsubscribe(topics...)
	}
}

func (c *Client) Subscribe(topics ...string) {
	c.client.AddConsumeTopics(topics...)
}

func (c *Client) Unsubscribe(topics ...string) {
	c.client.PurgeTopicsFromConsuming(topics...)
}

func (c *Client) Produce(msg ...types.Message) {
	if c.closing {
		return
	}

	for i := range msg {
		r := kgo.KeySliceRecord(msg[i].Key, msg[i].Value)
		r.Topic = msg[i].Topic
		// TODO: Headers!
		c.client.Produce(c.ctx, r, nil)
	}
}

func (c *Client) produceCallback(r *kgo.Record, err error) {
	if err != nil {
		c.logger.Error(c.facility, "Error producing message to '%s': %v",
			r.Topic, err)
	} else {
		c.logger.Trace(c.facility, "Message delivered to '%s'", r.Topic)
	}
}

func (c *Client) Close() {
	c.closing = true
	close(c.stop)
	c.cancel()

	// Some time to stop gracefully
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
	}

	c.cancel()
	c.client.CloseAllowingRebalance()
	c.logger.Info(c.facility, "Kafka client closed")
}
