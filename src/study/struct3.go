package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Card struct {
	ID    int     `json:"id,string"`
	Score float64 `json:"score,string"`
}

//    前端在传递来的json数据中可能会使用字符串类型的数字，
//   这个时候可以在结构体tag中添加string来告诉json包从对应字段解析相应的数据

func main() {
	str_json := `{"id":"123","score":"12.231"}`
	var card Card
	err := json.Unmarshal([]byte(str_json), &card)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(card)
	fmt.Println(reflect.TypeOf(card.ID))
}
