// Start runs a fixed-size worker pool.
//
// It starts workerCount goroutines that concurrently read tasks from inputCh,
// apply fn to each value, and send results to the returned output channel.
//
// The returned channel is closed when all workers exit. Workers exit when:
//   - ctx is cancelled, or
//   - inputCh is closed and fully drained.
package workerpool

import (
	"context"
	"sync"
)

func Start[T any, R any](ctx context.Context, workerCount int, inputCh <-chan T, fn func(e T) R) <-chan R {
	outputCh := make(chan R)
	wg := new(sync.WaitGroup)

	for range workerCount {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-inputCh:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case outputCh <- fn(val):
					}

				}

			}
		})
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}
