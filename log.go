package log

import (
	"context"
	"io"
	"os"
	"sync"
)

type level string

const (
	LevelInfo    level = "info"
	LevelError   level = "error"
	LevelUnknown level = "unknown"
)

var converters = make(map[string]func(interface{}) (interface{}, error))
var logger = New(os.Stderr)

type Logger struct {
	w io.Writer
	m sync.Mutex
}

// Error logs the passed data with LevelError level. The passed context is used
// to fetch request ID.
func (l *Logger) Error(ctx context.Context, data interface{}, restOfData ...interface{}) {
	l.log(LevelError, ctx, data, restOfData...)
}

// Info logs the passed data with LevelInfo level. The passed context is used
// to fetch request ID.
func (l *Logger) Info(ctx context.Context, data interface{}, restOfData ...interface{}) {
	l.log(LevelInfo, ctx, data, restOfData...)
}

func (l *Logger) log(level level, ctx context.Context, data interface{}, restOfData ...interface{}) {
	l.m.Lock()
	defer l.m.Unlock()
	l.w.Write(serialize(level, ctx, data, restOfData...))
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.m.Lock()
	defer l.m.Unlock()
	l.w = w
}

func New(w io.Writer) *Logger {
	return &Logger{w: w}
}

// Error logs the passed data with LevelError level. The passed context is used
// to fetch request ID.
func Error(ctx context.Context, data interface{}, restOfData ...interface{}) {
	logger.log(LevelError, ctx, data, restOfData...)
}

// Info logs the passed data with LevelInfo level. The passed context is used
// to fetch request ID.
func Info(ctx context.Context, data interface{}, restOfData ...interface{}) {
	logger.log(LevelInfo, ctx, data, restOfData...)
}

// GetWriter returns io.Writer that will log everything written to it under the
// specified level.
func GetWriter(level level) io.Writer {
	return &outputWriter{level, logger}
}

// SetOutput sets the output io.Writer for logging.
func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

// AddDataSerializationConverter adds a new serialization converter with the specified
// name. Serialization converters are used to create another serializable
// representation for logged data. For example if the data can't be serialized
// using the built in json.Marshal, you can supply a serialization converter
// to provide a representation that can be serialized. Data that were serialized
// from a serialization converter will be logged under "converted" key and then
// under the name of the serialization converter.
func AddDataSerializationConverter(name string, convert func(interface{}) (interface{}, error)) {
	converters[name] = convert
}
type outputWriter struct {
	level level
	l *Logger
}

func (w *outputWriter) Write(p []byte) (n int, err error) {
	w.l.log(w.level, context.TODO(), string(p))
	return len(p), nil
}