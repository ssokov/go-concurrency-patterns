package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ssokov/go-concurrency-patterns/internal/generator"
)

const (
	bufferSize  = 4
	workerCount = 2
)

func main() {
	var data []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*50)
	defer cancel()

	ch := generator.Generate(ctx, data, bufferSize)

	wg := new(sync.WaitGroup)

	for range workerCount {
		wg.Go(func() {

			for {
				select {
				case <-ctx.Done():
					return

				case val, ok := <-ch:
					if !ok {
						return
					}
					fmt.Println(val)
				}
			}
		})

	}

	wg.Wait()
}
