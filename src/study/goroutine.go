package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	pool *GoroutinePool
)

type GoroutinePool struct {
	c  chan struct{}
	wg *sync.WaitGroup
}

// 采用有缓冲channel实现,当channel满的时候阻塞
func NewGoroutinePool(maxSize int) *GoroutinePool {
	if maxSize <= 0 {
		panic("max size too small")
	}
	return &GoroutinePool{
		c:  make(chan struct{}, maxSize),
		wg: new(sync.WaitGroup),
	}
}

// add
func (g *GoroutinePool) Add(delta int) {
	g.wg.Add(delta)
	for i := 0; i < delta; i++ {
		g.c <- struct{}{}
	}

}

// done
func (g *GoroutinePool) Done() {
	<-g.c
	g.wg.Done()
}

// wait
func (g *GoroutinePool) Wait() {
	g.wg.Wait()
}

func main() {
	pool = NewGoroutinePool(100)
	done := make(chan struct{})
	pool.Add(5)
	for i := 0; i < 5; i++ {
		go run1()
	}
	go func() {
		pool.Wait()
		close(done)
	}()
	select {
	case <-done:
		fmt.Println("任务处理完成")
	case <-time.After(3 * time.Second):
		panic("处理任务超时")
	}
	fmt.Println("-------")
	pool.Add(7)
	for i := 0; i < 7; i++ {
		go run1()
	}
	go func() {
		pool.Wait()
		close(done)
	}()
	select {
	case <-done:
		fmt.Println("任务处理完成")
	case <-time.After(5 * time.Second):
		panic("处理任务超时")
	}
}

func run1() {
	fmt.Println(111)
	time.Sleep(4 * time.Second)
	fmt.Println(222)
	pool.Done()
}
