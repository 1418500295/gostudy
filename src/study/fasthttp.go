package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

func main() {
	req := &fasthttp.Request{}
	req.SetRequestURI("http://192.168.128.156:3333/boss/login")
	data := make(map[string]interface{})
	data["userName"] = "admin"
	data["password"] = "111"
	data["safeCode"] = "121"
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	req.SetBody(bytes)
	//req.Header.SetContentType()
	req.Header.SetMethod("POST") //必须用大写
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}
	err1 := client.Do(req, resp)
	if err1 != nil {
		fmt.Println(err1)
	}
	body := resp.Body()
	fmt.Println(string(body))
}
