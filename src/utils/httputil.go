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
	body, _ := res.Body()
	var jsonResp map[string]interface{}
	err1 := json.Unmarshal(body, &jsonResp)
	if err != nil {
		fmt.Println(err1)
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
