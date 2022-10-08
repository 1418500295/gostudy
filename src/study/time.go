package main

import (
	"fmt"
	"time"
)

func main() {
	//fmt.Println(timestamppb.New(time.Now()))
	////a := timestamppb.New(time.Now())
	//fmt.Println(timestamppb.Now())
	fmt.Println(time.Now().Month())
	for i := 0; i < 4; i++ {
		//定时器，设置两秒后执行
		t := time.NewTimer(2 * time.Second)
		<-t.C
		fmt.Println(i)
	}
}
