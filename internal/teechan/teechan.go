// Package teechan provides Teechan concurrency primitive
//
// Teechan duplicate values from a sungle input channel
// and sends their copies to mulitiple channels
//
// Each value from the input channel is forwarded to every our channel
// Idiomatically teechan should block, if even one out channel doesn't read
// it is like this, because this pattern tries to implement single channel
// logic using multiple, so it should be blocked until all out channels won't read
//
// Context is passed to exit blocking writing to out channel
package teechan

import "context"

func Teechan[T any](ctx context.Context, ch <-chan T, count int) []chan T {
	chs := make([]chan T, count)
	for i := range chs {
		chs[i] = make(chan T)
	}

	go func() {
		defer func() {
			for _, outCh := range chs {
				close(outCh)
			}
		}()

		for el := range ch {
			for _, ch := range chs {
				select {
				case <-ctx.Done():
					return
				case ch <- el:
				}
			}
		}
	}()

	return chs
}
