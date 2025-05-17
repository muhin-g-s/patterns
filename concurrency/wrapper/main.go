package main

import (
	"fmt"
	"sync"
	"time"
)

func wrapper(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		fmt.Println("Work before func")
		fn()
		time.Sleep(1 * time.Second)
		fmt.Println("Work after func")
	}()

}

func main() {
	var wg sync.WaitGroup

	wrapper(&wg, func() {
		time.Sleep(1 * time.Second)
		fmt.Println("heavy work")
	})

	wg.Wait()
}
