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
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	//刚开始先写入一条数据
	f, _ := os.OpenFile("./a.json", os.O_WRONLY, 0666)
	data := GetToken()
	reLi = append(reLi, data)
	bS, _ := json.Marshal(reLi)
	_, err = f.Write(bS)
	if err != nil {
		return
	}

	for i := 0; i < 10; i++ {
		//每次追加前先读取出来
		f1, _ := os.Open("./a.json")
		bs, err2 := io.ReadAll(f1)
		if err2 != nil {
			return
		}
		var rL []map[string]string
		err = json.Unmarshal(bs, &rL)

		data1 := GetToken()
		rL = append(rL, data1)
		//将新数据追加到list，再写入
		f2, _ := os.OpenFile("./a.json", os.O_WRONLY, 0666)
		bs1, _ := json.Marshal(rL)
		f2.Write(bs1)

	}
}




