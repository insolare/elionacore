package eliona

import (
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

type Logger struct{}

func NewLogger(level int) Logger {
	log.SetLevel(log.Level(level))

	return Logger{}
}

func (l Logger) Fatal(facility string, message string, params ...any) {
	log.Fatal(facility, message, params...)
}

func (l Logger) Trace(facility string, message string, params ...any) {
	log.Debug(facility, message, params...)
}

func (l Logger) Debug(facility string, message string, params ...any) {
	log.Fatal(facility, message, params...)
}

func (l Logger) Error(facility string, message string, params ...any) {
	log.Error(facility, message, params...)
}

func (l Logger) Warning(facility string, message string, params ...any) {
	log.Fatal(facility, message, params...)
}

func (l Logger) Info(facility string, message string, params ...any) {
	log.Fatal(facility, message, params...)
}
