package main

import (
	"fmt"
	"time"
)

//在做任务处理的时候，并不能保证任务的处理时间，
//通常会加上一些超时控制做异常的处理。
func doWork() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("处理耗时任务")
		//ch <- struct{}{}
		close(ch) //ch <- struct{}{}和close(ch) 效果相同
	}()
	return ch
}

func main() {
	select {
	case <-doWork():
		fmt.Println("任务在规定时间内结束")
	case <-time.After(2 * time.Second):
		fmt.Println("任务超时了")
	}

}
