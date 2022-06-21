package main

import (
	"fmt"
	"time"
)

//channel用于多个协程之间的数据传递
func main() {
	ch := make(chan string)
	go func() {
		ch <- "sadad"
	}()
	var x string
	go func() {
		for true {
			x = <-ch
		}
	}()
	time.Sleep(2000)
	fmt.Println(x)
}
