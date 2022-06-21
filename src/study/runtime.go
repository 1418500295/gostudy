package main

import (
	"fmt"
	"runtime"
)

func main() {

	fmt.Println(runtime.NumCPU())                     //获取系统cpu核数
	fmt.Println(runtime.GOMAXPROCS(runtime.NumCPU())) //设置最大的可用使用cpu核数
	runtime.Gosched()                                 //让当前线程让出 cpu 以让其它线程运行,
	// 它不会挂起当前线程，因此当前线程未来会继续执行
	runtime.GC()     //会让运行时系统进行一次强制性的垃圾收集
	runtime.Goexit() //退出当前 goroutine(但是defer语句会照常执行)
}
