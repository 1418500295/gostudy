package main

import (
	"fmt"
	"sync"
)

var sum1 = 0
var mutex1 sync.RWMutex

/**
读写锁：
适合读多写少的场景
1.RLock()，加读锁。某个goroutine加了读锁后，其它goroutine可以获取读锁，但是不能获取写锁
2.RUnlock()，释放读锁
3.加写锁。某个goroutine加了写锁后，其它goroutine不能获取读锁，也不能获取写锁
4.释放写锁。
读锁适合读多写少的场景
*/
func main() {

	for i := 0; i < 100; i++ {
		go func() {
			mutex1.RLocker()
			sum1 += 1
			defer mutex1.RUnlock()
		}()
	}
	fmt.Println(sum)
}
