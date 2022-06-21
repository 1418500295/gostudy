package main

import "fmt"

type Felines interface {
	feet()
}

type Land interface {
	work()
	Felines
}

type Chicken struct {
	name string
	age  int
}

func (c Chicken) work() {
	fmt.Println("c work")
}

func (c Chicken) feet() {
	fmt.Println("c feet")
}

func main() {
	var land Land = Chicken{"cc", 2}
	land.work()
	land.feet()

}
