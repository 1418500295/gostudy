package cases

import (
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"testing"
)

func Test_One(t *testing.T) {
	data := make(map[string]interface{})
	data["name"] = "daine"
	data["age"] = "26"
	//fmt.Println(utils.DoGetNoParams("http://localhost:8889/v1/getDemo1"))
	res, _ := HttpRequest.Get("http://localhost:8889/v1/getDemo1")
	var d map[string]interface{}
	err := res.Json(&d)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
}
