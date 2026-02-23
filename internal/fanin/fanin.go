// Package fanin provides fan-in concurrency primitives.
//
// FanIn merges multiple input channels into a single output channel.
//
// Each input channel is consumed concurrently in its own goroutine.
// Values received from any input channel are forwarded to the result
// channel.
package fanin

import (
	"context"
	"sync"
)

func FanIn[T any](ctx context.Context, buffer int, s ...<-chan T) <-chan T {
	result := make(chan T, buffer)

	wg := new(sync.WaitGroup)

	for _, ch := range s {

		wg.Go(func() {

			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-ch:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case result <- val:
					}

				}
			}

		})

	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}
