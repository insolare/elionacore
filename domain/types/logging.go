package types

type Logger interface {
	Fatal(string, string, ...any)
	Trace(string, string, ...any)
	Debug(string, string, ...any)
	Error(string, string, ...any)
	Warning(string, string, ...any)
	Info(string, string, ...any)
}

type NoopLogger struct{}

func (l NoopLogger) Fatal(string, string, ...any)   {}
func (l NoopLogger) Trace(string, string, ...any)   {}
func (l NoopLogger) Debug(string, string, ...any)   {}
func (l NoopLogger) Error(string, string, ...any)   {}
func (l NoopLogger) Warning(string, string, ...any) {}
func (l NoopLogger) Info(string, string, ...any)    {}
