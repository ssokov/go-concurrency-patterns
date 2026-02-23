// Package ratelimit provides rate limiting primitives.
//
// LeakyBuckerLimiter implements a classic leaky bucket algorithm.
// It allows up to N events per time period with a constant drain rate.
package ratelimit

import (
	"runtime"
	"time"
)

type LeakyBuckerLimiter struct {
	bucket   chan struct{}
	quit     chan struct{}
	interval time.Duration
}

func NewLeakyBuckerLimiter(limit int, period time.Duration) *LeakyBuckerLimiter {
	interval := period / time.Duration(limit)

	limiter := &LeakyBuckerLimiter{
		bucket:   make(chan struct{}, limit),
		quit:     make(chan struct{}),
		interval: interval,
	}

	go limiter.startLeaking()

	return limiter
}

func (l *LeakyBuckerLimiter) startLeaking() {

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

func (l *LeakyBuckerLimiter) Stop() {
	close(l.quit)
}

func (l *LeakyBuckerLimiter) Allow() bool {
	select {
	case l.bucket <- struct{}{}:
		return true
	default:
		return false
	}
}
