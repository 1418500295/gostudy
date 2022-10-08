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

//关闭 channel 一般是用来通知其他协程某个任务已经完成了
//
//关闭 channel 时应该注意以下准则：
//1)不要在读取端关闭 channel ，因为写入端无法知道 channel 是否已经关闭，往已关闭的 channel 写数据会 panic ；
//2)有多个写入端时，不要在写入端关闭 channle ，因为其他写入端无法知道 channel 是否已经关闭，关闭已经关闭的 channel 会发生 panic ；
//3)如果只有一个写入端，可以在这个写入端放心关闭 channel 。
//关闭 channel 粗暴一点的做法是随意关闭，如果产生了 panic 就用 recover 避免进程挂掉。稍好一点的方案是使用标准库的 sync 包来做关闭 channel 时的协程同步，不过使用起来也稍微复杂些
