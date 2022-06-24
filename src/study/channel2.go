package main

import (
	"fmt"
	"time"
)

//struct{}channel用户等待某事件先执行
func main() {
	channel := make(chan struct{})
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println(222)
		//close(channel)
		channel <- struct{}{} //或者使用close可以达到同样效果
	}()
	<-channel //阻塞下方程序执行，等待调用close后才执行
	fmt.Println(111)
}
