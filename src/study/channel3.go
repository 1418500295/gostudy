package main

import (
	"fmt"
)

func writeToChan(chanel chan int) {
	for i := 0; i < 4; i++ {
		chanel <- i
	}
	// 数据发送完毕，关闭通道
	close(chanel)

}

func main() {
	channel1 := make(chan int)
	go writeToChan(channel1)
	//for循环取完channel里的值后，因为通道close了，再次获取会拿到对应数据类型的零值
	//如果通道不close，for循环取完数据后就会阻塞报错
	for {
		v, ok := <-channel1
		if ok {
			fmt.Println(v)
			fmt.Println(ok)
		} else {
			fmt.Println("读取完成")
			fmt.Println(ok)
			break
		}
	}
}
