package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, result chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d start for %d\n", id, job)
		time.Sleep(time.Second)
		result <- job * 2
		fmt.Printf("Worker %d stop for %d\n", id, job)
	}
}

func main() {
	const SIZE_CHANNEL = 10

	jobs := make(chan int, SIZE_CHANNEL)
	results := make(chan int, SIZE_CHANNEL)

	const COUNT_WORKERS = 3

	for w := 1; w <= COUNT_WORKERS; w++ {
		go worker(w, jobs, results)
	}

	for i := 1; i <= SIZE_CHANNEL; i++ {
		jobs <- i
	}

	close(jobs)

	for i := 1; i <= SIZE_CHANNEL; i++ {
		fmt.Println("Result:", <-results)
	}
}
