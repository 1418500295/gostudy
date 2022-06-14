package cases

import (
	"fmt"
	"gostudy/src/utils"
	"testing"
)

const path = "C:\\Users\\AA\\go\\src\\gostudy\\src"

func Test_Two(t *testing.T) {
	resp := utils.DoPost(utils.GetApiUrl(path, "loginUri"),
		utils.GetTestData(path, "login.json", 0))
	fmt.Println(resp)
	//assert.Equal(t, -1, resp["code"])
	if resp["code"] != -2 {
		t.Error("失败")
	}
	
}

func Test_Three(t *testing.T) {
	fmt.Println(111)
}
