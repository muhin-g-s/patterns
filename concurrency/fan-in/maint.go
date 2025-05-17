package main

import (
	"fmt"
	"sync"
	"time"
)

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}

		wg.Done()
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	c1 := make(chan int)
	c2 := make(chan int)

	const (
		START  int = 0
		MIDDLE int = 5
		END    int = 10
	)

	go func() {
		for i := START; i < MIDDLE; i++ {
			c1 <- i
			time.Sleep(time.Second)
		}
		close(c1)
	}()

	go func() {
		for i := MIDDLE; i <= END; i++ {
			c2 <- i
			time.Sleep(time.Second)
		}
		close(c2)
	}()

	for n := range merge(c1, c2) {
		fmt.Println("Recived: ", n)
	}
}
