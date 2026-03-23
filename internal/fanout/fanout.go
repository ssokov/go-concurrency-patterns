// Package fanout provides fan-out concurrency primitives.
//
// FanOut distributes values from a single input channel
// across multiple output channels using a round-robin strategy.
//
// Each value from the input channel is forwarded to exactly
// one output channel. Distribution is performed sequentially
// in cyclic order.
package fanout

import (
	"context"
)

func FanOut[T any](ctx context.Context, inputCh <-chan T, size int) []chan T {

	if size <= 0 {
		return nil
	}

	outputCh := make([]chan T, size)

	for i := 0; i < len(outputCh); i++ {
		outputCh[i] = make(chan T)
	}

	go func() {
		defer func() {
			for _, c := range outputCh {
				close(c)
			}
		}()

		idx := 0

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
				case outputCh[idx] <- val:
				}

			}
			// Round Robin
			idx = (idx + 1) % size

		}

	}()

	return outputCh
}
