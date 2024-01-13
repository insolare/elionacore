package kafka

import "errors"

var (
	ErrNoSeeds = errors.New("no seed brokers provided")
)
