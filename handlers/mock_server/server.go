package main

import (
	"context"
	"log"
	"net/http"

	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/handlers/mock_server/di"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func run(router handlers.Router, logger gateways.Logger) error {
	m := http.NewServeMux()
	m.Handle("/", router)

	static := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	m.Handle("/static/", static)

	logger.Debugf(nil, "Open http://localhost:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", m); err != nil {
		return errors.Wrapf(err, "error while listening on port")
	}
	return nil
}

func main() {
	c, err := di.New()
	if err != nil {
		log.Fatalf("error while setting up a container: %+v", err)
	}
	if err := c.Provide(newZapLogger); err != nil {
		log.Fatalf("error while providing Logger: %+v", err)
	}
	if err := c.Invoke(run); err != nil {
		log.Fatalf("error while invoking app: %+v", err)
	}
}

func newZapLogger() (gateways.Logger, error) {
	// skip 1st caller that is zapLogger
	logger, err := zap.NewDevelopment(zap.AddCallerSkip(1))
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating a logger")
	}
	return &zapLogger{logger.Sugar()}, nil
}

type zapLogger struct {
	sugar *zap.SugaredLogger
}

func (l *zapLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

func (l *zapLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

func (l *zapLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

func (l *zapLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.sugar.Errorf(format, args)
}
