package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().AddDate(0, 0, 3))      //获取3天后时间
	fmt.Println(time.Now().Add(-5 * time.Minute)) //获取5分钟前时间
	fmt.Println(time.Now())
}
