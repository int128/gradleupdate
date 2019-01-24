package testing_logger

import (
	"context"
	"testing"

	"github.com/int128/gradleupdate/gateways/interfaces"
)

func New(t *testing.T) gateways.Logger {
	return &logger{t}
}

type logger struct {
	t *testing.T
}

func (l *logger) log(level string, format string, args ...interface{}) {
	l.t.Logf("["+level+"]"+format, args...)
}

func (l *logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.log("DEBUG", format, args...)
}

func (l *logger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

func (l *logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.log("WARN", format, args...)
}

func (l *logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}
