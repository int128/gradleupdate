package gatewaysTestDoubles

import (
	"context"
	"testing"
)

func NewLogger(t *testing.T) *Logger {
	return &Logger{t}
}

type Logger struct {
	t *testing.T
}

func (l *Logger) log(level string, format string, args ...interface{}) {
	l.t.Logf("["+level+"]"+format, args...)
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.log("DEBUG", format, args...)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

func (l *Logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.log("WARN", format, args...)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}
