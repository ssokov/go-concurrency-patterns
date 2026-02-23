package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ssokov/go-concurrency-patterns/internal/workerpool"
)

const (
	workerCount = 2
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ch := make(chan int, 100)

	for i := range 100 {
		ch <- i
	}
	close(ch)

	out := workerpool.Start(ctx, workerCount, ch, func(a int) int {
		time.Sleep(time.Millisecond * 100)
		return a * a
	})

	for val := range out {
		fmt.Println(val)
	}

}
