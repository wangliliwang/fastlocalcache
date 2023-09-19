package fastlocalcache

import "time"

type Clocker interface {
	CurrentTime() time.Time
}

type DefaultClock struct {
}

func (dc *DefaultClock) CurrentTime() time.Time {
	return time.Now()
}
