package main

import (
	"fmt"
	"sync"
)

func main() {
	var channel = make(chan int)
	var wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			channel <- i
		}
	}()
	go func() {
		defer wg.Done()
		for i := range channel {
			fmt.Println(i)
		}
	}()
	wg.Wait()
}
