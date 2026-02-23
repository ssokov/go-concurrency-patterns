// Generate starts a goroutine that emits elements of data to a channel.
//
// The returned channel is closed after all elements are sent
// or when the context is cancelled.
//
// The size parameter controls the channel buffer size.
// If size == 0, the channel is unbuffered and sends block
// until a receiver is ready.
package generator

import "context"

func Generate[T any](ctx context.Context, data []T, size int) <-chan T {

	result := make(chan T, size)

	go func() {
		defer close(result)

		for i := 0; i < len(data); i++ {
			select {
			case <-ctx.Done():
				return
			case result <- data[i]:
			}

		}

	}()

	return result
}
