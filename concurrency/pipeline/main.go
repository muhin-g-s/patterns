package main

import (
	"fmt"
	"time"
)

type StageFunc func(int) int

func gen(max int) <-chan int {
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

func stage(in <-chan int, fn StageFunc) <-chan int {
	out := make(chan int)

	go func() {
		for n := range in {
			out <- fn(n)
		}
		close(out)
	}()

	return out
}

func main() {
	multiply := func(n int) int { return n * 2 }
	square := func(n int) int { return n * n }
	add := func(n int) int { return n + 10 }

	pipeline := stage(stage(stage(gen(5), multiply), square), add)

	fmt.Println("Pipeline results:")
	for n := range pipeline {
		fmt.Println(n)
	}
}
