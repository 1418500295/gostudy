package main

import (
	"encoding/json"
	"fmt"
)

// Person /**
type Person struct {
	Name   string  `json:"name"`             //指定序列化时用小写name
	Age    int     `json:"-"`                //指定序列化忽略该字段
	Weight float64 `json:"weight,omitempty"` //omitempty:当struct中该字段无值时，序列化后的结果要忽略时，需添加omitempty
}

func main() {
	person := Person{"", 12, 21.3}
	byteS, _ := json.Marshal(person)
	fmt.Println(string(byteS))
}
