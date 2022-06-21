package main

import "fmt"

type Phone interface {
	call()
	sendMsg()
}

type Xiaomi struct {
	name  string
	price int
}

type Huawei struct {
	name  string
	price int
}

/**
1.只要有某个方法的实现使用了指针接受者，那给包含了这个方法的interface变量赋值的时候要使用指针。
比如上面的Dog类型要赋值给Animal，必须使用指针，因为Dog实现speak方法用了指针接受者
2.如果全部方法都使用的是值接受者，那给interface变量赋值的时候用值或者指针都可以
*/
func (huawei Huawei) call() {
	fmt.Printf("%s有打电话功能", huawei.name)
}

func (huawei *Huawei) sendMsg() {
	fmt.Printf("%s可以发短信", huawei.name)
}

func main() {
	//mate30 := Huawei{"Mate30", 1000}
	//mate30.call()
	//mate30.sendMsg()
	var phone Phone
	phone = &Huawei{"Mate40", 2000}
	phone.call()
	phone.sendMsg()
}
