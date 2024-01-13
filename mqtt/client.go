package mqtt

import (
	"context"

	paho3 "github.com/eclipse/paho.mqtt.golang"
	"github.com/insolare/elionacore/domain/types"
)

type Client struct {
	client     paho3.Client
	publishQoS byte
}

type ClientConfig struct {
	Broker     any
	PublishQoS byte
}

var globalClient *Client

func Init() {}

func NewClient(ctx context.Context, cfg Client, fn HandleFunc) (*Client, error) {
	opts := paho3.NewClientOptions()
	opts.SetDefaultPublishHandler(makeMqttHandler(fn))

	return nil, nil
}

func (c *Client) Subscribe(topics ...string) {
	for i := range topics {
		c.client.Subscribe(topics[i], 0, nil)
	}
}

func (c *Client) Unsubscribe(topics ...string) {
	for i := range topics {
		c.client.Unsubscribe(topics[i])
	}
}

func (c *Client) Produce(messages ...types.Message) {
	for i := range messages {
		c.client.Publish(
			messages[i].Topic,
			c.publishQoS,
			false,
			messages[i].Value,
		)
	}
}
