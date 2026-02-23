package main

import (
	"context"

	"github.com/ssokov/go-concurrency-patterns/internal/fanin"
	"github.com/ssokov/go-concurrency-patterns/internal/fanout"
	"github.com/ssokov/go-concurrency-patterns/internal/generator"
)

func worker[T any, R any](ctx context.Context, in <-chan T, fn func(T) R) <-chan R {
	out := make(chan R)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- fn(val):
				}
			}
		}
	}()

	return out
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	input := generator.Generate(ctx, data, 1)

	// fan-out on N channels
	N := 3
	outs := fanout.FanOut(ctx, input, N)

	// starting N workers
	workerOuts := make([]<-chan int, 0, N)

	for _, ch := range outs {
		workerOuts = append(workerOuts,
			worker(ctx, ch, func(v int) int {
				return v * v
			}),
		)
	}

	// fan-in
	result := fanin.FanIn(ctx, 0, workerOuts...)

	for v := range result {
		println(v)
	}
}
