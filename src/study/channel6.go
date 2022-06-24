package main

import (
	"fmt"
	"time"
)

//假设有4个 goroutine，编号为1，2，3，4。每秒钟会有一个 goroutine 打印出它自己的编号。
//现在让你写一个程序，要求输出的编号总是按照1，2，3，4这样的顺序打印

type token struct{}

func main() {
	num := 4
	var chs []chan token
	// 4 个work
	for i := 0; i < num; i++ {
		chs = append(chs, make(chan token))
	}
	for j := 0; j < num; j++ {
		go worker(j, chs[j], chs[(j+1)%num])
	}
	// 先把令牌交给第一个
	chs[0] <- struct{}{}
	select {}
}

func worker(id int, ch chan token, next chan token) {
	for {
		// 对应work 取得令牌
		token := <-ch
		fmt.Println(id + 1)
		time.Sleep(1 * time.Second)
		// 传递给下一个
		next <- token
	}
}
