package types

// Message is a message sent to or received from Eliona.
// Depending on transport (Kafka, MQTT), fields may be used differently internaly
type Message struct {
	Topic   string
	Key     []byte
	Value   []byte
	Headers any // Need to make a wrapper
}
