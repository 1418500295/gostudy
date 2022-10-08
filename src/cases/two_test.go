package main

import (
	"fmt"
	"testing"
)

func Test_Two(t *testing.T) {
	fmt.Println(222)
	//resp := utils.DoPost(utils.GetApiUrl("loginUri"),
	//	utils.GetTestData("login.json", 0))
	//fmt.Println(resp)
	////assert.Equal(t, -1, resp["code"])
	//if resp["code"] != 0 {
	//	t.Error("失败")
	//}

}

func Test_Three(t *testing.T) {
	fmt.Println(333)
}

//
//func Test_Three(t *testing.T) {
//	fmt.Println(utils.GetTestData(path1, "login.json", 0))
//
//}
