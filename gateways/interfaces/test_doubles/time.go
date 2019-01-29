package gateways

import "time"

type TimeService struct {
	NowValue time.Time
}

func (s *TimeService) Now() time.Time {
	return s.NowValue
}
