package main

import (
	"fmt"
	"time"
)

func generateNumbers(max int) chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for i := 1; i <= max; i++ {
			out <- i

			time.Sleep(time.Second)
		}
	}()

	return out
}

func main() {
	for i := range generateNumbers(10) {
		fmt.Println("number ", i)
	}
}
