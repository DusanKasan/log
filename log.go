package log

import (
	"context"
	"io"
	"os"
)

type level string

const (
	LevelInfo    level = "info"
	LevelError   level = "error"
	LevelUnknown level = "unknown"
)


var writer io.Writer = os.Stderr
var converters = make(map[string]func(interface{}) (interface{}, error))

type Logger interface {
	// Error logs the passed data with LevelError level.
	Error(data interface{}, restOfData ...interface{})
	// Info logs the passed data with LevelInfo level.
	Info(data interface{}, restOfData ...interface{})
}

// Error logs the passed data with LevelError level. The passed context is used
// to fetch request ID.
func Error(ctx context.Context, data interface{}, restOfData ...interface{}) {
	writeWithLevel(LevelError, ctx, data, restOfData...)
}

// Info logs the passed data with LevelInfo level. The passed context is used
// to fetch request ID.
func Info(ctx context.Context, data interface{}, restOfData ...interface{}) {
	writeWithLevel(LevelInfo, ctx, data, restOfData...)
}

// WithContext returns a Logger instance that has access to the passed context.
func WithContext(ctx context.Context) Logger {
	return &contextLogger{ctx}
}

// GetWriter returns io.Writer that will log everything written to it under the
// specified level.
func GetWriter(level level) io.Writer {
	return &outputWriter{level}
}

// SetOutput sets the output io.Writer for logging.
func SetOutput(w io.Writer) {
	writer = w
}

// AddDataSerializationConverter adds a new serialization converter with the sepcified
// name. Serialization converters are used to create another serializable
// representation for logged data. For example if the data can't be serialized
// using the built in json.Marshal, you can supply a serialization converter
// to provide a representation that can be serialized. Data that were serialized
// from a serialization converter will be logged under "converted" key and then
// under the name of the serialization converter.
func AddDataSerializationConverter(name string, convert func(interface{}) (interface{}, error)) {
	converters[name] = convert
}

type contextLogger struct {
	ctx context.Context
}

// Error logs the passed data with LevelError level.
func (i *contextLogger) Error(data interface{}, restOfData ...interface{}) {
	Error(i.ctx, data, restOfData...)
}

// Info logs the passed data with LevelInfo level.
func (i *contextLogger) Info(data interface{}, restOfData ...interface{}) {
	Info(i.ctx, data, restOfData...)
}

type outputWriter struct {
	level level
}

func (w *outputWriter) Write(p []byte) (n int, err error) {
	writeWithLevel(w.level, context.TODO(), string(p))
	return len(p), nil
}
