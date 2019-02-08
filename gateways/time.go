package gateways

import (
	"time"

	"go.uber.org/dig"
)

type Time struct {
	dig.In
}

func (s *Time) Now() time.Time {
	return time.Now()
}
