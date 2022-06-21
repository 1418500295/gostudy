package main

import (
	"fmt"
	"sync"
)

var sum = 0
var mutex sync.Mutex

//多个 goroutine对共享变量同时执行写操作，并发是不安全的，结果和预期不符。
func main() {

	for i := 0; i < 100; i++ {
		go func() {
			mutex.Lock()
			sum += 1
			defer mutex.Unlock()
		}()
	}
	fmt.Println(sum)
}
