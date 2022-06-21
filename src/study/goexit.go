package main

import (
	"fmt"
	"runtime"
	"time"
)

/**
runtime.Goexit() //退出当前 goroutine(但是defer语句会照常执行)
*/

func main() {
	go func() {
		fmt.Println("开始")
		one()
		fmt.Println("结束")
	}()
	time.Sleep(3 * time.Second)

}

func one() {
	defer fmt.Println("我是defer函数")
	runtime.Goexit()
	fmt.Println(111)
}
