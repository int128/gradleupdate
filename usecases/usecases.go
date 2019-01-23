package usecases

import "time"

type TimeProvider func() time.Time

func (p TimeProvider) Now() time.Time {
	if p == nil {
		return time.Now()
	}
	return p()
}
