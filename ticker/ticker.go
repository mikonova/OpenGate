package ticker

import (
	"sync"
	"time"
)

type PrecisionTicker struct {
	offsetTime   time.Duration
	intervalTime time.Duration
	sigChan      chan int
	Tick         chan time.Time
	mut          sync.RWMutex
}

type Offset struct {
	SecondOffset     int64
	NanosecondOffset int64
}

func NewPTicker(offset, interval time.Duration) *PrecisionTicker {
	pt := PrecisionTicker{
		offsetTime:   offset,
		intervalTime: interval,
		sigChan:      make(chan int, 1),
		Tick:         make(chan time.Time),
	}
	go func() {
		time.After(pt.offsetTime)
		for {
			if n := <-pt.sigChan; n == -1 {
				break
			}
			pt.mut.Lock()
			pt.Tick <- time.Now()
			pt.mut.Unlock()
			time.Sleep(pt.intervalTime)
		}
	}()
	return &pt
}

func (t *PrecisionTicker) Close() {
	defer close(t.sigChan)
	defer close(t.Tick)
	t.sigChan <- -1
}
