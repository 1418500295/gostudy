package main

import (
	"fmt"
	"time"
)

//select监听多路channel，哪个先满足就执行那哪个，执行完就退出
func main() {
	cha1 := make(chan int)
	cha2 := make(chan int)
	cha3 := make(chan int)
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println(111)
		cha1 <- 1
		close(cha1)
	}()
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println(222)
		cha2 <- 2
		close(cha2)
	}()
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println(333)
		cha3 <- 3
		close(cha3)
	}()

	select {
	case v := <-cha1:
		fmt.Println(v)
	case v := <-cha2:
		fmt.Println(v)
	case v := <-cha3:
		fmt.Println(v)
	}

}
