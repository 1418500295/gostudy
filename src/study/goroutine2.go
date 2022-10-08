package main

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

//获取单个协程占用内存
func main() {

	var c chan int
	var wg sync.WaitGroup
	const gotun = 1e4

	mem := func() uint64 {
		runtime.GC()
		var memStat runtime.MemStats
		runtime.ReadMemStats(&memStat)
		return memStat.Sys
	}
	noop := func() {
		wg.Done()
		<-c
	}
	wg.Add(gotun)
	before := mem()
	for i := 0; i < gotun; i++ {
		go noop()
	}
	wg.Wait()
	after := mem()
	fmt.Println(float64(after-before) / gotun / 1000)
	var a = 5
	fmt.Println(reflect.TypeOf(a))

}
