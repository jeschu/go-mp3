package eta

import (
	"fmt"
	"sync"
	"time"
)

type Eta struct {
	mu    sync.Mutex
	start time.Time
	count uint64
	total uint64
}

func NewEta(total uint64) Eta {
	return Eta{start: time.Now(), count: 0, total: total}
}

func (eta *Eta) SetCount(count uint64) { eta.count = count }
func (eta *Eta) IncCount() {
	eta.mu.Lock()
	eta.count++
	eta.mu.Unlock()
}
func (eta *Eta) SetTotal(total uint64) { eta.total = total }

func (eta *Eta) String() string {
	now := time.Now()
	rt := now.Sub(eta.start)
	fCount := float64(eta.count)
	fTotal := float64(eta.total)
	fLeft := float64(eta.total - eta.count)
	done := fCount / fTotal * 100.0
	perStep := float64(rt) / float64(eta.count)
	left := time.Duration(fLeft * perStep)
	return fmt.Sprintf(
		"%s (runtime: %s) %d/%d %6.2f%% => %s (%s)",
		eta.start.Format("15:04:05"),
		rt.Truncate(time.Second),
		eta.count,
		eta.total,
		done,
		(left + time.Second).Truncate(time.Second),
		now.Add(left).Format("15:04:05"),
	)
}

func (eta *Eta) ETA() time.Time {
	now := time.Now()
	return now.Add(time.Duration(float64(eta.start.Sub(now)) / (float64(eta.count) / float64(eta.total))))
}
