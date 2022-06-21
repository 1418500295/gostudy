package main

import "fmt"

/**
一个struct类型可以实现多个interface。比如猫这个类型，既是猫科动物，也是哺乳动物。
猫科动物可以是一个interface，哺乳动物可以是另一个interface，猫这个struct类型可以实现猫科动物和哺乳动物这2个interface里的方法
*/
type Animal interface {
	speck()
}
type Mammal interface {
	born()
}

type Cat struct {
	name string
	age  int
}

func (cat Cat) speck() {
	fmt.Println("cat speck")
}

func (cat Cat) born() {
	fmt.Println("cat born")
}

func main() {
	//var animal Animal = Cat{"rick", 1}
	//animal.speck()
	//var mammal Mammal = Cat{"rick", 2}
	//mammal.born()

	cat := Cat{"rich", 3}
	var animal Animal = cat
	animal.speck()
	var mammal Mammal = cat
	mammal.born()
}
