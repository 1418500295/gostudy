package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetTestData(projectPath string, fileName string, caseIndex int) map[string]interface{} {
	//path, _ := os.Getwd()
	byteData, err := ioutil.ReadFile(projectPath + "/testdata/" + fileName)
	if err != nil {
		fmt.Println(err)
	}
	var jsonData []map[string]interface{}
	err1 := json.Unmarshal(byteData, &jsonData)
	if err1 != nil {
		fmt.Println(err1)
	}
	return jsonData[caseIndex]
}

func GetApiUrl(projectPath string, urlName string) string {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("捕获到异常: ", err)
		}
	}()

	//path, err := os.Getwd()
	//if err != nil {
	//	fmt.Println(err)
	//}
	files, err1 := os.Open(projectPath + "/host.properties")
	defer files.Close()
	if err1 != nil {
		fmt.Println(err1)
	}
	bytesStr, _ := ioutil.ReadAll(files)
	configStr := strings.Split(string(bytesStr), "\n")
	var endUrl string
	var host string
	for _, i := range configStr {
		iSlice := strings.Split(strings.ReplaceAll(i, "\r", ""), "=")
		if "host" == strings.Trim(iSlice[0], " ") {
			host = strings.Trim(iSlice[1], " ")
		}
		if urlName == strings.Trim(iSlice[0], " ") {
			url := strings.Trim(iSlice[1], " ")
			if !strings.HasPrefix(url, "/") {
				endUrl = host + "/" + fmt.Sprintf("%v", url)
			} else if strings.HasPrefix(url, "/") {
				endUrl = host + fmt.Sprintf("%v", url)
			} else {
				panic("请求地址格式不正确")
			}
		}
	}
	return endUrl
}
