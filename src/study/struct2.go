package main

import (
	"encoding/json"
	"fmt"
)

// Person1 嵌套结构体 /**
type Person1 struct {
	Age    int
	Height float64
	Weight float64
}

type Person2 struct {
	Age int
	Sex int
}

type User struct {
	Name     string
	Hobby    []string
	Person1  Person1                    `json:"person1"` //如果需要序列化旧的person为嵌套json，据需要打tag为别名
	*Person2 `json:"person2,omitempty"` //如果需要序列化后的数据无该结构体，需要标记为指针类型
}

func main() {
	user := User{Name: "bls", Hobby: []string{"篮球", "足球"}, Person1: Person1{
		1, 2.4, 4.5,
	}}
	byteS, _ := json.Marshal(user)
	fmt.Println(string(byteS))
}
