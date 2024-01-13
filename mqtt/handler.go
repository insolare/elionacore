package mqtt

import (
	"dev.azure.com/eliona/libs-go/_git/core/domain/types"
	paho3 "github.com/eclipse/paho.mqtt.golang"
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
