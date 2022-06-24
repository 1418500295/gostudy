package utils

import (
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
)

func DoPost(url string, testData map[string]interface{}) map[string]interface{} {
	res, err := HttpRequest.Post(url, testData)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Close()
	defer fmt.Println("响应时间：" + fmt.Sprintf("%v", res.Time()))
	var jsonResp map[string]interface{}
	err1 := res.Json(&jsonResp)
	if err1 != nil {
		fmt.Println(err)
	}
	//判断map中的值类型，因为map解析会将int值自动转为float64，并将值(一般为interface类型)转为string或int
	for k, v := range jsonResp {
		switch v.(type) {
		case float64:
			jsonResp[k] = int(v.(float64))
		case string:
			jsonResp[k] = v.(string)
		}

	}
	defer res.Close()
	return jsonResp
}

func DoJsonPost(url string, testData map[string]interface{}) map[string]interface{} {
	res, err := HttpRequest.JSON().
		Post(url, testData)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := res.Body()
	defer res.Close()
	defer fmt.Println("响应时间：" + fmt.Sprintf("%v", res.Time()))
	var jsonRes map[string]interface{}
	err1 := json.Unmarshal(body, &jsonRes)
	if err1 != nil {
		fmt.Println(err1)
	}
	for k, v := range jsonRes {
		switch v.(type) {
		case float64:
			jsonRes[k] = int(v.(float64))
		case string:
			jsonRes[k] = v.(string)
		}
	}
	return jsonRes

}

func DoGet(url string, testData map[string]interface{}) map[string]interface{} {
	res, err := HttpRequest.Get(url, testData)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := res.Body()
	defer res.Close()
	defer fmt.Println("响应时间: " + fmt.Sprintf("%v", res.Time()))
	var jsonRes map[string]interface{}
	err1 := json.Unmarshal(body, &jsonRes)
	if err1 != nil {
		fmt.Println(err1)
	}
	for k, v := range jsonRes {
		switch v.(type) {
		case float64:
			jsonRes[k] = int(v.(float64))
		case string:
			jsonRes[k] = v.(string)
		}
	}
	return jsonRes

}
func DoGetNoParams(url string) map[string]interface{} {
	res, err := HttpRequest.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := res.Body()
	defer res.Close()
	defer fmt.Println("响应时间：" + fmt.Sprintf("%v", res.Time()))
	var jsonRes map[string]interface{}
	err1 := json.Unmarshal(body, &jsonRes)
	if err1 != nil {
		fmt.Println(err1)
	}
	for k, v := range jsonRes {
		switch v.(type) {
		case float64:
			jsonRes[k] = int(v.(float64))
		case string:
			jsonRes[k] = v.(string)
		}
	}
	return jsonRes
}
