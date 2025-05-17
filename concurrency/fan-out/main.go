package main

import (
	"fmt"
	"sync"
	"time"
)

func split(ch <-chan int, n int) []chan int {
	cs := make([]chan int, n)

	for i := 0; i < n; i++ {
		cs[i] = make(chan int)
	}

	distributeToChannels := func(ch <-chan int, cs []chan int) {
		defer func(cs []chan int) {
			for _, c := range cs {
				close(c)
			}
		}(cs)

		for {
			for _, c := range cs {
				select {
				case val, ok := <-ch:
					if !ok {
						return
					}
					c <- val
				}
			}
		}
	}

	go distributeToChannels(ch, cs)
	return cs
}

func producer(tasks chan<- int, numTasks int) {
	defer close(tasks)
	for i := 1; i <= numTasks; i++ {
		time.Sleep(time.Second)
		fmt.Printf("Producer: sending task %d\n", i)
		tasks <- i
	}
}

func consumer(id int, wg *sync.WaitGroup, jobs <-chan int, result chan<- int) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Consumer %d: start for %d\n", id, job)
		time.Sleep(time.Second)
		result <- job * 2
		fmt.Printf("Consumer %d: stop for %d\n", id, job)
	}
}

func main() {
	const NUM_TASKS = 10

	tasks := make(chan int, NUM_TASKS)
	result := make(chan int, NUM_TASKS)
	var wg sync.WaitGroup

	go producer(tasks, NUM_TASKS)

	cs := split(tasks, NUM_TASKS)

	for i, c := range cs {
		wg.Add(1)
		go consumer(i, &wg, c, result)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	var results []int
	for r := range result {
		results = append(results, r)
	}

	fmt.Println("\nFinal results:")
	for i, r := range results {
		fmt.Printf("Result %d: %d\n", i+1, r)
	}
}
