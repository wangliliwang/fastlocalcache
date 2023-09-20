package fastlocalcache

import "time"

type Clocker interface {
	CurrentTime() time.Time
	CurrentTimestamp() uint64
}

type DefaultClock struct {
}

func (dc *DefaultClock) CurrentTime() time.Time {
	return time.Now()
}

func (dc *DefaultClock) CurrentTimestamp() uint64 {
	return uint64(dc.CurrentTime().Unix())
}
