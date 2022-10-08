package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randInt(min, max int) int {
	//rand.New(rand.NewSource(time.Now().Unix()))
	rand.Seed(time.Now().Unix()) //设置种子
	return rand.Intn(max-min) + min
}
func main() {
	//设置种子
	//rand.Seed(time.Now().Unix())
	//fmt.Println(randInt(10,30))
	fmt.Println(randInt(10, 20))
}
