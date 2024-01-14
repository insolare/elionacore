package mqtt

import (
	paho3 "github.com/eclipse/paho.mqtt.golang"
	"github.com/insolare/elionacore/domain/types"
)

type HandleFunc func(m types.Message)

func makeMqttHandler(handler HandleFunc) paho3.MessageHandler {
	return func(_ paho3.Client, m paho3.Message) {
		msg := types.Message{
			Topic:   m.Topic(),
			Key:     nil,
			Value:   m.Payload(),
			Headers: nil,
		}

		handler(msg)
	}
}
