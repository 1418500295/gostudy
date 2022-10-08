package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	wsLiveUri       = "ws://chat.mliveplus.com/"
	wsType          = "login"
	sendDataType    = "say"
	sendContent     = "斗皇强者，恐怖如斯！😂😂😂😂"
	content2        = "仙路尽头谁为峰，一遇无始道成空!🤙🤙🤙🤙🤙"
	content3        = "天不生我李淳刚，剑道万古如长夜!😭😭😭😭"
	content4        = "吾儿王腾有大帝之资～！😡😡😡😡"
	content5        = "遮天，谁以弹指而遮天？天地不仁，何以万物争相残？蚍蜉撼树可献命，是因其心蛇吞仙！🌚🌚🌚🌚🌚🌚"
	content6        = "黑袍老者目中精光一闪，旋即倒吸一口凉气！👽👽👽"
	roomId          = 2
	platForm        = 0
	loginType       = 1
	loginUrl        = "https://api.mliveplus.com/api/user/login"
	orderUrl        = "https://api.mliveplus.com/api/charge/orderList"
	withDrawsUrl    = "https://api.mliveplus.com/api/withdraw/withdraws"
	getSpendListUrl = "https://api.mliveplus.com/api/order/getSpendList"
	getLiveUrl      = "https://api.mliveplus.com/webapi/live/getLivePageData"
	getRecordUrl    = "https://sport.sun8tv.com/api/task/getRecord"
	register1       = "https://api.mliveplus.com/webapi/user/regist"
	register2       = "https://api.mliveplus.com/webapi/user/register"
	path            = "/Users/eden/go/src/gostudy/src/ws"
)

var (
	//go:embed data.json
	f embed.FS
)
var (
	wg              = sync.WaitGroup{}
	num             int64
	okNumLive       int64 = 0
	channelLive           = make(chan int64)
	resTimeLiveList []int
	lockLive              = sync.Mutex{}
	iLive           int64 = 0
	doneLive1             = make(chan struct{})
	//resTime     int64
	oneLive        int64
	twoLive        int64
	threeLive      int64
	firstLive      int64
	secondLive     int64
	thirdLive      int64
	connectLiveNum int64 = 0
	token          string
	uid            string
	nickName       string
	hash           string
	accountList    []map[string]string
	userData       []map[string]string
	account        string
	pwd            string
	avatar         string
	phone          string

	hashChan  = make(chan string)
	tokenChan = make(chan map[string]interface{})
)

func maxRespTimeLive() int {
	max := resTimeLiveList[0]
	for _, i := range resTimeLiveList {
		if i > max {
			max = i
		}
	}
	return max
}
func minRespTimeLive() int {
	min := resTimeLiveList[0]
	for _, i := range resTimeLiveList {
		if i < min {
			min = i
		}
	}
	return min
}
func sumResTimeLive() int {
	sum := 0
	for _, i := range resTimeLiveList {
		sum += i
	}
	return sum
}
func fiftyRespTime() int {
	sort.Ints(resTimeLiveList)
	resSize := 0.5
	return resTimeLiveList[int(float64(len(resTimeLiveList))*resSize)-1]
}
func ninetyRespTime() int {
	sort.Ints(resTimeLiveList)
	resSize := 0.9
	return resTimeLiveList[int(float64(len(resTimeLiveList))*resSize)-1]
}

//注册
func exeRegister() {
	done := make(chan struct{})
	wg.Add(int(num))
	for i := 0; i < int(num); i++ {
		go registerSetup()
		go register(hashChan)
		go func() {
			data := <-channelLive
			resTimeLiveList = append(resTimeLiveList, int(data))
		}()
	}
	wg.Wait()
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	select {
	case <-done:
		fmt.Println("***任务处理完成***")
	case <-time.After(time.Duration(30) * time.Second):
		fmt.Println("***任务处理超时***")
	}
}

//登陆
func loginExe() {
	done := make(chan struct{})
	wg.Add(int(num))
	for i := 0; i < int(num); i++ {
		go login(i)
		go func() {
			data := <-channelLive
			resTimeLiveList = append(resTimeLiveList, int(data))
		}()
	}
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	select {
	case <-done:
		fmt.Println("***任务处理完成***")
	case <-time.After(time.Duration(30) * time.Second):
		fmt.Println("***任务处理超时***")
	}
}
func runLive(funName func(chan map[string]interface{})) {
	done := make(chan struct{})
	wg.Add(int(num))
	for i := 0; i < int(num); i++ {
		go loginSetUp(i)
		go funName(tokenChan)
		go func() {
			data := <-channelLive
			resTimeLiveList = append(resTimeLiveList, int(data))
		}()
	}
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	select {
	case <-done:
		fmt.Println("***任务处理完成***")
	case <-time.After(time.Duration(30) * time.Second):
		fmt.Println("***任务处理超时***")
	}
}

func main() {
	byteData1, err := f.ReadFile("data.json")
	if err != nil {
		fmt.Println(err)
	}
	err1 := json.Unmarshal(byteData1, &userData)
	if err1 != nil {
		fmt.Println(err1)
	}
	apiList := []string{"(1)orderList", "(2)withDraws", "(3)getSpendList", "(4)websocket", "(5)register", "(6)login"}
	fmt.Println("接口列表：", apiList)
	var apiNum int64
	fmt.Println("请选择并发接口序号(如：1): ")
	_, _ = fmt.Scan(&apiNum)
	fmt.Println("请输入并发数：")
	_, _ = fmt.Scan(&num)
	sTime := time.Now().UnixNano() / 1e6
	switch apiNum {
	case 1:
		runLive(orderList)
	case 2:
		runLive(withDraws)
	case 3:
		runLive(getSpendList)
	case 4:
		runLive(liveSend)
	case 5:
		exeRegister()
	case 6:
		loginExe()
	default:
		fmt.Println("请输入正确的接口序号")
	}
	eTime := time.Now().UnixNano() / 1e6
	fmt.Println(fmt.Sprintf("总耗时：%v秒", float64(eTime-sTime)/1000))
	fmt.Println("总并发数：", num)
	fmt.Println("总请求数: ", num)
	fmt.Printf("成功的数量：%v \n", okNumLive)
	fmt.Printf("失败的数量：%v \n", num-okNumLive)
	fmt.Printf("最小响应时间：%.3f微秒 ≈ %v秒 \n", float64(minRespTimeLive()), float64(minRespTimeLive())/1e6)
	fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTimeLive())/1e6)
	fmt.Println("50%用户响应时间: " + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(fiftyRespTime()), float64(fiftyRespTime())/1e6))
	fmt.Println("90%用户响应时间: " + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(ninetyRespTime()), float64(ninetyRespTime())/1e6))
	fmt.Println(fmt.Sprintf("平均响应时间：%.3f秒", (float64(sumResTimeLive())/float64(len(resTimeLiveList)))/1e6))
	fmt.Println(fmt.Sprintf("QPS: %.3f", float64(num)/(float64(sumResTimeLive())/float64(len(resTimeLiveList))/1e6)))
	defer runtime.GC()
}

func registerSetup() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	req := HttpRequest.NewRequest()
	payLoad := make(map[string]string)
	rand.Seed(time.Now().UnixNano())
	payLoad["mobile"] = fmt.Sprintf("151%v3%v", rand.Intn(999-100)+100, rand.Intn(9999-1000)+1000)
	payLoad["smscode"] = "999999"
	res, _ := req.JSON().Post(register1, payLoad)
	body, _ := res.Body()
	var resMap map[string]interface{}
	err := json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	hashChan <- gjson.ParseBytes(body).Map()["data"].Map()["mobile_hash"].String()
	phone = payLoad["mobile"]
	fmt.Println(res.StatusCode())
	defer res.Close()

}
func register(ch chan string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	req := HttpRequest.NewRequest()
	rand.Seed(time.Now().UnixNano())
	payLoad := make(map[string]interface{})
	ranS := rand.Intn(999-110) + 110
	ranS1 := rand.Intn(999-110) + 110
	payLoad["account"] = fmt.Sprintf("den%v%v1", ranS, ranS1)
	payLoad["mobile_hash"] = <-ch
	payLoad["invite_code"] = ""
	payLoad["platform"] = 0
	payLoad["sex"] = 0
	payLoad["pwd"] = fmt.Sprintf("den%v%v1", ranS, ranS1)
	res, _ := req.JSON().Post(register2, payLoad)
	body, _ := res.Body()
	resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
	channelLive <- resTime * 1000
	var resMap map[string]interface{}
	err := json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
		atomic.AddInt64(&okNumLive, 1)
	} else {
		log.Println("响应异常：", string(body))
	}
	log.Println("响应码：", res.StatusCode())
	log.Println("接口返回：", string(body))
	//fmt.Println(string(body))
	//account = gjson.ParseBytes(body).Map()["data"].Map()["account"].String()
	//pwd = gjson.ParseBytes(body).Map()["data"].Map()["account"].String()
}

func login(i int) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	req := HttpRequest.NewRequest()
	reqUrl := loginUrl
	payLoad := make(map[string]interface{})
	payLoad["account"] = userData[i]["account"]
	payLoad["login_type"] = loginType
	payLoad["platform"] = platForm
	payLoad["pwd"] = userData[i]["pwd"]
	res, _ := req.JSON().Post(reqUrl, payLoad)
	body, _ := res.Body()
	resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
	channelLive <- resTime * 1000
	var resMap map[string]interface{}
	err := json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
		atomic.AddInt64(&okNumLive, 1)
	} else {
		log.Println("响应异常：", string(body))
	}
	log.Println("响应码：", res.StatusCode())
	log.Println("接口返回：", string(body))
	defer wg.Done()
}
func loginSetUp(i int) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	req := HttpRequest.NewRequest()
	reqUrl := loginUrl
	payLoad := make(map[string]interface{})
	payLoad["account"] = userData[i]["account"]
	payLoad["login_type"] = loginType
	payLoad["platform"] = platForm
	payLoad["pwd"] = userData[i]["pwd"]
	res, _ := req.JSON().Post(reqUrl, payLoad)
	defer res.Close()
	body, _ := res.Body()
	chMap := make(map[string]interface{})
	chMap["token"] = gjson.ParseBytes(body).Map()["data"].Map()["token"].String()
	chMap["id"] = gjson.ParseBytes(body).Map()["data"].Map()["id"].String()
	chMap["nick_name"] = gjson.ParseBytes(body).Map()["data"].Map()["nick_name"].String()
	chMap["avatar"] = gjson.ParseBytes(body).Map()["data"].Map()["avatar"].String()
	tokenChan <- chMap
}
func orderList(ch chan map[string]interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	req := HttpRequest.NewRequest()
	data := make(map[string]interface{})
	data["page"] = 1
	data["size"] = 10
	chMap := <-ch
	headers := make(map[string]string)
	headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
	fmt.Println(headers)
	res, _ := req.SetHeaders(headers).JSON().Post(orderUrl, data)
	resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
	channelLive <- resTime * 1000
	body, _ := res.Body()
	defer res.Close()
	var resMap map[string]interface{}
	err1 := json.Unmarshal(body, &resMap)
	if err1 != nil {
		fmt.Println(err1)
	}
	if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
		atomic.AddInt64(&okNumLive, 1)
	}
	fmt.Println(reflect.TypeOf(resMap["status"]))
	log.Println("响应码：", res.StatusCode())
	log.Println("接口返回：", string(body))
	defer wg.Done()
}

func withDraws(ch chan map[string]interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	req := HttpRequest.NewRequest()
	data := make(map[string]interface{})
	data["end-time"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
	data["start_time"] = strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
	data["page"] = 1
	data["size"] = 10
	data["status"] = 0
	chMap := <-ch
	headers := make(map[string]string)
	headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
	res, _ := req.SetHeaders(headers).JSON().Post(withDrawsUrl, data)
	resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
	channelLive <- resTime * 1000
	body, _ := res.Body()
	defer res.Close()
	var resMap map[string]interface{}
	err1 := json.Unmarshal(body, &resMap)
	if err1 != nil {
		log.Println("解析返回数据异常: ", err1)
	}
	if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
		atomic.AddInt64(&okNumLive, 1)
	}
	log.Println("响应码：", res.StatusCode())
	log.Println("接口返回：", string(body))
	defer wg.Done()
}

func getSpendList(ch chan map[string]interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	req := HttpRequest.NewRequest()
	data := make(map[string]interface{})
	data["etimestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
	data["stimestamp"] = strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
	data["page"] = 1
	data["size"] = 10
	data["coin_type"] = 0
	data["spend_type"] = 0
	chMap := <-ch
	headers := make(map[string]string)
	headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
	res, _ := req.SetHeaders(headers).JSON().Post(getSpendListUrl, data)
	resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
	channelLive <- resTime * 1000
	body, _ := res.Body()
	res.Close()
	var resMap map[string]interface{}
	err1 := json.Unmarshal(body, &resMap)
	if err1 != nil {
		log.Println("解析返回数据异常: ", err1)
	}
	if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
		atomic.AddInt64(&okNumLive, 1)
	}
	log.Println("响应码：", res.StatusCode())
	log.Println("接口返回：", string(body))
	defer wg.Done()

}

func liveSend(ch chan map[string]interface{}) {
	defer func() {
		err5 := recover()
		if err5 != nil {
			log.Println("捕获的异常：", err5)
		}
	}()
	chMap := <-ch
	liveData := make(map[string]interface{})
	liveData["type"] = wsType
	liveData["client_name"] = chMap["nick_name"]
	liveData["room_id"] = roomId
	userId, _ := strconv.Atoi(chMap["id"].(string))
	liveData["user_id"] = userId
	liveData["avatar"] = chMap["avatar"]
	liveData["token"] = <-tokenChan
	liveData["platform"] = platForm
	bs, _ := json.Marshal(&liveData)
	ws, res, err1 := websocket.DefaultDialer.Dial(wsLiveUri, nil)
	if err1 != nil {
		log.Println(fmt.Sprintf("建立连接异常：%v", err1))
	}
	if res != nil {
		log.Println(`连接成功: `, res)
		atomic.AddInt64(&connectLiveNum, 1)
	}
	log.Println("已成功建立的连接数：", connectLiveNum)
	err2 := ws.WriteMessage(websocket.BinaryMessage, bs)
	log.Println(`请求进入直播间 ：`, string(bs))
	if err2 != nil {
		log.Println("进入直播间异常：", err2)
	}
	_, recv1, err9 := ws.ReadMessage()
	if err9 != nil {
		log.Println("进入直播间异常", err9)
	}
	log.Println(`进入直播间成功：`, string(recv1))
	contentlist := []string{sendContent, content2, content3, content4, content5, content6}
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	sendD := make(map[string]interface{})
	sendD["type"] = sendDataType
	sendD["content"] = contentlist[source.Intn(5)]
	sendBs, _ := json.Marshal(&sendD)
	sTime := time.Now().UnixNano() / 1e3
	err8 := ws.WriteMessage(websocket.BinaryMessage, sendBs)
	log.Println(`发送聊天信息：`, string(sendBs))
	if err8 != nil {
		log.Println("发送聊天信息异常：", err2)
	}
	_, recv, err3 := ws.ReadMessage()
	if err3 != nil {
		log.Println("接收数据异常", err3)
	}
	var resMsg map[string]interface{}
	err4 := json.Unmarshal(recv, &resMsg)
	if err4 != nil {
		log.Println("解析返回数据异常：", err4)
	}
	log.Println(`服务端返回：`, string(recv))
	eTime := time.Now().UnixNano() / 1e3
	if string(recv) != "" {
		atomic.AddInt64(&okNumLive, 1)
	} else {
		fmt.Println("接收服务端返回数据异常：", string(recv))
	}
	channelLive <- eTime - sTime
	defer ws.Close()
	defer wg.Done()
}
