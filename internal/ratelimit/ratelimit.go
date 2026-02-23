// Package ratelimit provides rate limiting primitives.
//
// leakyBuckerLimiter implements a classic leaky bucket algorithm.
// It allows up to N events per time period with a constant drain rate.
//
// A background goroutine is started internally to leak tokens at a fixed interval.
// Callers should invoke Stop() when the limiter is no longer needed
// to terminate the background goroutine. It is recommended to use:
//
//	rl := ratelimit.NewleakyBuckerLimiter(...)
//	defer rl.Stop() 
//
// to ensure proper resource cleanup.
package ratelimit

import (
	"runtime"
	"sync"
	"time"
)

type leakyBuckerLimiter struct {
	bucket   chan struct{}
	quit     chan struct{}
	interval time.Duration
	once     sync.Once
}

func NewLeakyBuckerLimiter(limit int, period time.Duration) *leakyBuckerLimiter {
	interval := period / time.Duration(limit)

	limiter := &leakyBuckerLimiter{
		bucket:   make(chan struct{}, limit),
		quit:     make(chan struct{}),
		interval: interval,
	}

	go limiter.startLeaking()

	return limiter
}

func (l *leakyBuckerLimiter) startLeaking() {

	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case <-l.bucket:
			default:
				runtime.Gosched()
			}
		case <-l.quit:
			return
		}
	}

}

func (l *leakyBuckerLimiter) Stop() {
	l.once.Do(func() {
		close(l.quit)
	})

}

func (l *leakyBuckerLimiter) Allow() bool {
	select {
	case l.bucket <- struct{}{}:
		return true
	default:
		return false
	}
}
