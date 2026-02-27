package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/ssokov/go-concurrency-patterns/internal/generator"
	"github.com/ssokov/go-concurrency-patterns/internal/teechan"
)

type Teechan[T any] func(<-chan T, int) []chan T

func main() {
	ctx := context.Background()

	outCount := 2
	fmt.Println("out channels count:", outCount)
	ch := generator.Generate(ctx, []string{"Alex", "Bob", "Nataly"}, 0)

	chs := teechan.Teechan(ctx, ch, 2)

	wg := sync.WaitGroup{}
	for id, outCh := range chs {
		wg.Add(1)
		go func(c chan string, id int) {
			defer wg.Done()
			for i := range c {
				fmt.Printf("recieved %v from channel with id %d\n", i, id)
			}
		}(outCh, id)
	}
	wg.Wait()
}
