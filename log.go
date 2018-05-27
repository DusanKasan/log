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
	Error(data interface{}, restOfData ...interface{})
	Info(data interface{}, restOfData ...interface{})
}

func Error(ctx context.Context, data interface{}, restOfData ...interface{}) {
	writeWithLevel(LevelError, ctx, data, restOfData...)
}

func Info(ctx context.Context, data interface{}, restOfData ...interface{}) {
	writeWithLevel(LevelInfo, ctx, data, restOfData...)
}

func WithContext(ctx context.Context) *contextLogger {
	return &contextLogger{ctx}
}

func GetWriter(level level) io.Writer {
	return &outputWriter{level}
}

func SetWriter(w io.Writer) {
	writer = w
}

func AddSerializationConverter(name string, convert func(interface{}) (interface{}, error)) {
	converters[name] = convert
}

type contextLogger struct {
	ctx context.Context
}

func (i *contextLogger) Error(data interface{}, restOfData ...interface{}) {
	Error(i.ctx, data, restOfData...)
}

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
