package main

import (
	"fmt"
	"gostudy/src/utils"
	"testing"
)

func Test_Two(t *testing.T) {
	resp := utils.DoPost(utils.GetApiUrl("loginUri"),
		utils.GetTestData("login.json", 0))
	fmt.Println(resp)
	//assert.Equal(t, -1, resp["code"])
	if resp["code"] != 0 {
		t.Error("失败")
	}

}

//
//func Test_Three(t *testing.T) {
//	fmt.Println(utils.GetTestData(path1, "login.json", 0))
//
//}
