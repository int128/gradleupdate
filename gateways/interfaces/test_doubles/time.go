package gateways

import "time"

type FixedTime struct {
	NowValue time.Time
}

func (s *FixedTime) Now() time.Time {
	return s.NowValue
}
