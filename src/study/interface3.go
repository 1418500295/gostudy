package main

import "fmt"

/**
多个struct类型可以实现同一个interface：多个类型都有共同的方法(行为)
*/
type House interface {
	live()
}

type Dog struct {
	name string
	age  int
}

type Pig struct {
	name string
	age  int
}

func (dog Dog) live() {
	fmt.Println("dog live")

}

func (pip *Pig) live() {
	fmt.Println("pig live")
}

func main() {
	var house House = Dog{"aa", 1}
	house.live()

	house = &Pig{"bb", 22}
	house.live()
}
