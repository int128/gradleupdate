package gateways

import (
	"time"

	"go.uber.org/dig"
)

type TimeService struct {
	dig.In
}

func (s *TimeService) Now() time.Time {
	return time.Now()
}
