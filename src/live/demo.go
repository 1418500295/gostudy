package main

import (
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	register1Url = "https://ycapi.mliveplus.com/api/user/regist"
	register2Url = "https://ycapi.mliveplus.com/api/user/register"
	loginUrl2    = "https://ycapi.mliveplus.com/api/user/login"
	path2        = "/Users/eden/go/src/gostudy/src/live/data.json"
)

var (
	account111 string
	pwd111     string
	phone111   string
	accList    []map[string]string
	userD      []map[string]string
	g          = sync.WaitGroup{}
	l          sync.Mutex
)

func main() {
	//byteData, err := ioutil.ReadFile(path2)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err1 := json.Unmarshal(byteData, &userD)
	//if err1 != nil {
	//	fmt.Println(err1)
	//}

	//for i := 0; i < 10; i++ {
	//	register222()
	//}

	//for i := 0; i < 3; i++ {
	//	dengLu()
	//}
	byteData, err := ioutil.ReadFile(path2)
	if err != nil {
		fmt.Println(err)
	}
	var s []map[string]string
	err1 := json.Unmarshal(byteData, &s)
	if err1 != nil {
		fmt.Println(err1)
	}
	s = s[4836:]
	//s = append(s, accList...)
	//
	f1, _ := os.OpenFile(path2, os.O_TRUNC|os.O_WRONLY, 0666)
	bS, _ := json.Marshal(s)
	_, _ = f1.Write(bS)

	//err2 := json.Unmarshal(byteData1, &userD)
	//if err2 != nil {
	//	fmt.Println(err2)
	//}
	//register222()
	//headers := make(map[string]string)
	//headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", "hj+TzZwPBclbxLfgEf2uIJDwSAH0dlFi", "27224")
	//req := HttpRequest.NewRequest()
	//rand.Seed(time.Now().UnixNano())
	//card := rand.Intn((1000 - 1) + 1)
	//card2 := rand.Intn((1000 - 1) + 1)
	//payLoad := make(map[string]interface{})
	//payLoad["bank"] = "工商银行"
	//payLoad["card"] = fmt.Sprintf("%v%v", card, card2)
	//payLoad["id_card"] = fmt.Sprintf("%v%v", card, card2)
	//payLoad["name"] = "eden"
	//payLoad["smscode"] = "999999"
	//res, err2 := req.SetHeaders(headers).JSON().Post("http://ycapi.mliveplus.com/"+"api/user/addBank", payLoad)
	//if err2 != nil {
	//	log.Println("请求异常：", err2)
	//}
	//log.Println("响应码：", res.StatusCode())
	//log.Println("请求Url: ", 1)
	//body, _ := res.Body()
	//log.Println("接口返回：", string(body))

}

var j int64 = 0
var li []string

func dengLu() {

	req := HttpRequest.NewRequest()
	mP := register222()
	payLoad := make(map[string]interface{})
	payLoad["account"] = mP["account"]
	payLoad["login_type"] = 1
	payLoad["platform"] = 0
	payLoad["pwd"] = mP["pwd"]
	payLoad["source_id"] = 0
	res, _ := req.JSON().Post(loginUrl2, payLoad)
	body, _ := res.Body()
	res.Close()
	fmt.Println(string(body))
	if gjson.ParseBytes(body).Map()["status"].String() == "0" {
		//atomic.AddInt64(&j, 1)
		//l.Lock()
		//li = append(li, gjson.ParseBytes(body).Map()["data"].Map()["account"].String())
		//l.Unlock()
		dataMap := make(map[string]string)
		dataMap["account"] = account111
		dataMap["pwd"] = pwd111
		dataMap["phone"] = phone111
		dataMap["token"] = gjson.ParseBytes(body).Map()["data"].Map()["token"].String()
		dataMap["id"] = gjson.ParseBytes(body).Map()["data"].Map()["id"].String()
		dataMap["nick_name"] = gjson.ParseBytes(body).Map()["data"].Map()["nick_name"].String()
		dataMap["avatar"] = gjson.ParseBytes(body).Map()["data"].Map()["avatar"].String()
		accList = append(accList, dataMap)
	}

}

func registerSetup111() string {

	var err error
	req := HttpRequest.NewRequest()
	payLoad := make(map[string]string)
	rand.Seed(time.Now().UnixNano())
	payLoad["mobile"] = fmt.Sprintf("185%v3%v", rand.Intn(999-100)+100, rand.Intn(9999-1000)+1000)
	payLoad["smscode"] = "999999"
	res, err := req.JSON().Post(register1Url, payLoad)
	if err != nil {
		fmt.Println("请求异常：", err)
	}
	fmt.Println(register1Url)
	body, _ := res.Body()
	var resMap map[string]interface{}
	err = json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println("解析返回数据异常：", err)
	}
	fmt.Println(string(body))
	hash := gjson.ParseBytes(body).Map()["data"].Map()["mobile_hash"].String()
	phone111 = payLoad["mobile"]
	res.Close()
	return hash
}

func register222() map[string]string {
	var err error
	req := HttpRequest.NewRequest()
	rand.Seed(time.Now().UnixNano())
	payLoad := make(map[string]interface{})
	ranS := rand.Intn(999-110) + 110
	ranS1 := rand.Intn(999-110) + 110
	payLoad["account"] = fmt.Sprintf("aa%v%v2", ranS, ranS1)
	payLoad["mobile_hash"] = registerSetup111()
	payLoad["invite_code"] = ""
	payLoad["sex"] = 1
	payLoad["platform"] = 0
	payLoad["pwd"] = fmt.Sprintf("aa%v%v2", ranS, ranS1)
	//payLoad["reg_source_id"] = 10002
	res, err := req.JSON().Post(register2Url, payLoad)
	fmt.Println(register2Url)
	if err != nil {
		fmt.Println("请求异常：", err)
	}
	body, _ := res.Body()
	res.Close()
	//channelLive <- resTime * 1000
	var resMap map[string]interface{}
	err = json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println("解析返回异常：", err)
	}
	fmt.Println(string(body))
	account111 = gjson.ParseBytes(body).Map()["data"].Map()["account"].String()
	pwd111 = gjson.ParseBytes(body).Map()["data"].Map()["account"].String()
	zh := make(map[string]string)
	zh["account"] = account111
	zh["pwd"] = pwd111
	return zh
}
