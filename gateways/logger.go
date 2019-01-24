package gateways

import (
	"context"

	"go.uber.org/dig"
	"google.golang.org/appengine/log"
)

// AELogger provides a logger using appengine/log package.
type AELogger struct {
	dig.In
}

func (l *AELogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	log.Debugf(ctx, format, args...)
}

func (l *AELogger) Infof(ctx context.Context, format string, args ...interface{}) {
	log.Infof(ctx, format, args...)
}

func (l *AELogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	log.Warningf(ctx, format, args...)
}

func (l *AELogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(ctx, format, args...)
}
