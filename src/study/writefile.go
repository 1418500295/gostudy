package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//追加后覆盖写入文件
func main() {
	path := "/login.json"
	ByS, _ := ioutil.ReadFile(path)
	var ds []map[string]interface{}
	err := json.Unmarshal(ByS, &ds)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ds)
	data := make(map[string]interface{})
	data["name"] = "asd"
	data["age"] = 12
	ds = append(ds, data)
	file, _ := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0666)
	dataByte, _ := json.Marshal(ds)
	_, _ = file.Write(dataByte)
	defer file.Close()
	fmt.Println(ds)
}
