package main

import (
	"fmt"
	"time"

	"github.com/ssokov/go-concurrency-patterns/internal/ratelimit"
)

const (
	requestsPerInterval = 2
	intervalDuration    = time.Second
)

func main() {
	// 10 requst per 1 second
	rl := ratelimit.NewLeakyBuckerLimiter(requestsPerInterval, intervalDuration)

	for range 100 {
		if rl.Allow() {
			fmt.Println("Allowed")
		} else {
			fmt.Println("Not allowed")
		}
		time.Sleep(time.Millisecond * 100)
	}

}
