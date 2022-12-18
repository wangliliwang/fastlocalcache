package fastcache_lru

import "fmt"

type stats struct {
	GetCalls uint64
	SetCalls uint64
	DelCalls uint64
	Misses   uint64
}

func (s *stats) String() string {
	return fmt.Sprintf(`GetCalls: %+v
SetCalls: %+v
DelCalls: %+v
Misses:   %+v
MissRate: %.2f%%
`, s.GetCalls, s.SetCalls, s.DelCalls, s.Misses, float64(s.Misses)/float64(s.GetCalls)*100)
}
