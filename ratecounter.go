package ratecounter

import (
	"sync"
	"time"
)

type RateCounter struct {
	*sync.RWMutex
	counts   map[string]int
	total    int
	interval time.Duration
	add      chan string
	birth    time.Time
	stop 	 chan bool
}

func NewRateCounter(interval time.Duration) *RateCounter {
	return &RateCounter{
		RWMutex:  &sync.RWMutex{},
		counts:   make(map[string]int),
		add:      make(chan string),
		birth:    time.Now(),
		interval: interval,
		stop:     make(chan bool),
	}
}

func (r *RateCounter) Inc(tag string) {
	r.add <- tag
}

func (r *RateCounter) Start() {
	dec := make(chan string)
	go func() {
		for {
			select {
			case tag := <-r.add:
				r.Lock()
				interval := r.interval
				r.counts[tag]++
				r.total++
				r.Unlock()
				go func() {
					//<-time.After(interval)
					time.Sleep(interval)
					dec <- tag
				}()
			case tag := <-dec:
				r.Lock()
				r.counts[tag]--
				r.total--
				r.Unlock()
			case s := <-r.stop:
				if s {
					return
				}
			}
		}
	}()
}

func (r *RateCounter) Stop() {
	r.stop <- true
}

func (r *RateCounter) Count(tag string) int {
	r.RLock()
	defer r.RUnlock()
	return r.counts[tag]
}

func (r *RateCounter) CountAll() int {
	r.RLock()
	defer r.RUnlock()
	return r.total
}

func (r *RateCounter) Age() time.Duration {
	r.RLock()
	defer r.RUnlock()
	return time.Since(r.birth)
}

