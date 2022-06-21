package main

import (
	"fmt"
	"runtime"
)

func main() {
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println(1111)
		}
	}()

	for i := 0; i < 10; i++ {
		runtime.Gosched() //让出时间片，先让别的协议执行，它执行完，再回来执行此协程
		fmt.Println(2222)
	}
}
