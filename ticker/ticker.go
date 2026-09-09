package ticker

import (
	"sync"
	"time"
)

type PrecisionTicker struct {
	startTime    time.Time
	intervalTime time.Duration
	sig          int
	Tick         chan time.Time
	mut          sync.RWMutex
}

type Offset struct {
	SecondOffset     int64
	NanosecondOffset int64
}

func NewPTicker(ntpTime time.Time, interval time.Duration) *PrecisionTicker {
	pt := PrecisionTicker{
		startTime:    ntpTime,
		intervalTime: interval,
		Tick:         make(chan time.Time, 1),
	}
	go func() {
		for {
			pt.mut.Lock()
			if pt.sig == -1 {
				pt.mut.Unlock()
				break
			}
			pt.mut.Unlock()
			if len(pt.Tick) == 1 {
				_ = <-pt.Tick
			}
			pt.Tick <- ntpTime
			elapsed := time.Since(ntpTime)
			ntpTime = ntpTime.Add(elapsed)
			sleepTime := pt.intervalTime - elapsed
			time.Sleep(sleepTime)
		}
	}()
	return &pt
}

func (pt *PrecisionTicker) Close() {
	defer close(pt.Tick)
	pt.mut.RLock()
	pt.sig = -1
	pt.mut.RUnlock()
}
