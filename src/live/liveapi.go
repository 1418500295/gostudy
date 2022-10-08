package main

import (
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
	"log"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	platForm        = 0
	loginType       = 1
	account         = "carl001"
	pwd             = "carl001"
	loginUrl        = "https://sport.sun8tv.com/api/user/login"
	orderUrl        = "https://sport.sun8tv.com/api/charge/orderList"
	withDrawsUrl    = "https://sport.sun8tv.com/api/withdraw/withdraws"
	getSpendListUrl = "https://sport.sun8tv.com/api/order/getSpendList"
)

var (
	//创建计数器
	wg = sync.WaitGroup{}
	//num      int64 = 5 //设置并发数量
	okNum    int64 = 0 //初始化请求成功的数量
	timeList []int     //响应时间
	channel        = make(chan int64)
	done           = make(chan struct{})
	lock           = sync.Mutex{}
	iNum     int64 = 0
	one      int64
	two      int64
	three    int64
	first    int64
	second   int64
	third    int64
	token    string
	uid      int
)

type Response struct {
	Resp map[string]gjson.Result `json:"msg"`
}

//获取时间戳
func SetTs() string {
	return strconv.FormatInt(time.Now().Unix()*1000, 10)
}

func sumRespTime() int {
	sum := 0
	for _, index := range timeList {
		sum = sum + index
	}
	return sum
}

func maxRespTime() int {
	max := timeList[0]
	for _, index := range timeList {
		if index > max {
			max = index
		}
	}
	return max
}
func minRespTime() int {
	min := timeList[0]
	for _, index := range timeList {
		if index < min {
			min = index
		}
	}
	return min
}

func fiftyRespTime() int {
	sort.Ints(timeList)
	resSize := 0.5
	return timeList[int(float64(len(timeList))*resSize)-1]
}
func ninetyRespTime() int {
	sort.Ints(timeList)
	resSize := 0.9
	return timeList[int(float64(len(timeList))*resSize)-1]
}
func printTime(useTime int64, count int64) {
	fmt.Println("并发数: ", count)
	fmt.Println("请求数: ", iNum)
	fmt.Println("成功的数量：", okNum)
	fmt.Printf("失败的数量：%v \n", iNum-okNum)
	fmt.Println(fmt.Sprintf("耗时：%v秒", float64(useTime)/1000))
	fmt.Println("50%用户响应时间：" + fmt.Sprintf("%.3f秒", float64(fiftyRespTime())/1000))
	fmt.Println("90%用户响应时间：" + fmt.Sprintf("%.3f秒", float64(ninetyRespTime())/1000))
	fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTime())/1000)
	fmt.Printf("最小响应时间：%.3f毫秒 \n", float64(minRespTime()))
	fmt.Printf("平均响应时间是:%.3f秒 \n", float64(sumRespTime())/float64(iNum)/1000)
	fmt.Printf("QPS：%.3f \n", float64(count)/(float64(sumRespTime())/float64(iNum)/1000))
}
func login() {
	req := HttpRequest.NewRequest()
	reqUrl := loginUrl
	payLoad := make(map[string]interface{})
	payLoad["account"] = account
	payLoad["login_type"] = loginType
	payLoad["platform"] = platForm
	payLoad["pwd"] = pwd
	res, _ := req.JSON().Post(reqUrl, payLoad)
	body, _ := res.Body()
	var resMap map[string]interface{}
	err := json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println(err)
	}
	token = gjson.ParseBytes(body).Map()["data"].Map()["token"].String()
	uid = int(gjson.ParseBytes(body).Map()["data"].Map()["id"].Int())
	//nickName = gjson.ParseBytes(body).Map()["data"].Map()["nick_name"].String()
	//fmt.Println(string(body))
}
func init() {
	login()
}
func main() {
	//one1()
	var apiList = []string{"(1)orderList", "(2)withDraws", "(3)getSpendList", "(4)Bet", "(5)TransactionHistory"}
	fmt.Println("接口列表：", apiList)
	var apiName string
	fmt.Println("请输入要请求的接口编号如：1")
	_, _ = fmt.Scan(&apiName)
	fmt.Println("第一轮压测用户数:")
	_, _ = fmt.Scan(&one)
	fmt.Println("第一轮运行时长(秒):")
	_, _ = fmt.Scan(&first)
	fmt.Println("第二轮压测用户数：")
	_, _ = fmt.Scan(&two)
	fmt.Println("第二轮运行时长(秒):")
	_, _ = fmt.Scan(&second)
	fmt.Println("第三轮压测用户数：")
	_, _ = fmt.Scan(&three)
	fmt.Println("第三轮运行时长(秒):")
	_, _ = fmt.Scan(&third)
	startTime := time.Now().UnixNano() / 1e6
	fmt.Printf("开始时间：%v \n", startTime)
	do(apiName)
	endTime := time.Now().UnixNano() / 1e6
	fmt.Printf("结束时间：%v \n", endTime)
	fmt.Println("总并发数:", one+two+three)
	fmt.Println("总请求数: ", iNum)
	fmt.Println("成功的数量: ", okNum)
	fmt.Printf("失败的数量: %v \n", iNum-okNum)
	fmt.Printf("总耗时：%.3f 秒 \n", float64(endTime-startTime)/1000-(5+5))
	fmt.Println("50%用户响应时间: " + fmt.Sprintf("%.3f秒", float64(fiftyRespTime())/1000))
	fmt.Println("90%用户响应时间: " + fmt.Sprintf("%.3f秒", float64(ninetyRespTime())/1000))
	fmt.Printf("最大响应时间: %.3f秒 \n", float64(maxRespTime())/1000)
	fmt.Printf("最小响应时间: %.3f毫秒 \n", float64(minRespTime()))
	fmt.Printf("平均响应时间是:%.3f秒 \n", float64(sumRespTime())/float64(len(timeList))/1000)
	fmt.Printf("QPS: %.3f \n", float64(one+two+three)/(float64(sumRespTime())/float64(len(timeList))/1000))
	runtime.GC()
	//_, _ = fmt.Scanf("h")
}

//, one int64, first int64, two int64, second int64, three int64, third int64
func do(name string) {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	switch name {
	case "1":
		fmt.Println(fmt.Sprintf("***首轮并发协程为%v***", one))
		firSt := time.Now().UnixNano() / 1e6
		wg.Add(int(one))
		for i := 0; i < int(one); i++ {
			go OrderList(first)
		}
		wg.Wait()
		firEnd := time.Now().UnixNano() / 1e6
		fmt.Println("***第一轮压测结果***")
		printTime(firEnd-firSt, one)
		fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", two))
		<-time.After(10 * time.Second)
		seSt := time.Now().UnixNano() / 1e6
		wg.Add(int(two))
		for i := 0; i < int(two); i++ {
			go OrderList(second)
		}
		wg.Wait()
		seEnd := time.Now().UnixNano() / 1e6
		fmt.Println("***第二轮压测结果***")
		printTime(seEnd-seSt, two)
		fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", three))
		<-time.After(10 * time.Second)
		wg.Add(int(three))
		for i := 0; i < int(three); i++ {
			go OrderList(third)
		}
		wg.Wait()
	case "2":
		fmt.Println(fmt.Sprintf("***首轮并发协程为%v***", one))
		firSt := time.Now().UnixNano() / 1e6
		wg.Add(int(one))
		for i := 0; i < int(one); i++ {
			go withDraws(first)
		}
		wg.Wait()
		firEd := time.Now().UnixNano() / 1e6
		fmt.Println("***第一轮压测结果***")
		printTime(firEd-firSt, one)
		fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", two))
		<-time.After(10 * time.Second)
		seSt := time.Now().UnixNano() / 1e6
		wg.Add(int(two))
		for i := 0; i < int(two); i++ {
			go withDraws(second)
		}
		wg.Wait()
		seEd := time.Now().UnixNano() / 1e6
		fmt.Println("***第二轮压测结果***")
		printTime(seEd-seSt, two)
		fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", three))
		<-time.After(10 * time.Second)
		wg.Add(int(three))
		for i := 0; i < int(three); i++ {
			go withDraws(third)
		}
		wg.Wait()
	case "3":
		fmt.Println(fmt.Sprintf("***首轮并发协程为%v***", one))
		firSt := time.Now().UnixNano() / 1e6
		wg.Add(int(one))
		for i := 0; i < int(one); i++ {
			go getSpendList(first)
		}
		wg.Wait()
		firEd := time.Now().UnixNano() / 1e6
		fmt.Println("***第一轮压测结果***")
		printTime(firEd-firSt, one)
		fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", two))
		<-time.After(10 * time.Second)
		seSt := time.Now().UnixNano() / 1e6
		wg.Add(int(two))
		for i := 0; i < int(two); i++ {
			go getSpendList(second)
		}
		wg.Wait()
		seEd := time.Now().UnixNano() / 1e6
		fmt.Println("***第二轮压测结果***")
		printTime(seEd-seSt, two)
		fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", three))
		<-time.After(10 * time.Second)
		wg.Add(int(three))
		for i := 0; i < int(three); i++ {
			go getSpendList(third)
		}
		wg.Wait()
	case "4":
		fmt.Println(fmt.Sprintf("***首轮并发协程为%v***", one))
		firSt := time.Now().UnixNano() / 1e6
		wg.Add(int(one))
		for i := 0; i < int(one); i++ {
			//go Bet(first)
		}
		wg.Wait()
		firEd := time.Now().UnixNano() / 1e6
		fmt.Println("***第一轮结束后压测结果***")
		printTime(firEd-firSt, one)
		fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", two))
		<-time.After(10 * time.Second)
		seSt := time.Now().UnixNano() / 1e6
		wg.Add(int(two))
		for i := 0; i < int(two); i++ {
			//go Bet(second)
		}
		wg.Wait()
		seEd := time.Now().UnixNano() / 1e6
		fmt.Println("***第二轮结束后压测结果***")
		printTime(seEd-seSt, two)
		fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", three))
		<-time.After(10 * time.Second)
		wg.Add(int(three))
		for i := 0; i < int(three); i++ {
			//go Bet(third)
		}
		wg.Wait()
	case "5":
		fmt.Println(fmt.Sprintf("***首轮并发协程为%v***", one))
		firSt := time.Now().UnixNano() / 1e6
		wg.Add(int(one))
		for i := 0; i < int(one); i++ {
			//go TransactionHistory(first)
		}
		wg.Wait()
		firEd := time.Now().UnixNano() / 1e6
		fmt.Println("***第一轮结束后压测结果***")
		printTime(firEd-firSt, one)
		fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", two))
		<-time.After(10 * time.Second)
		seSt := time.Now().UnixNano() / 1e6
		wg.Add(int(two))
		for i := 0; i < int(two); i++ {
			//go TransactionHistory(second)
		}
		wg.Wait()
		seEd := time.Now().UnixNano() / 1e6
		fmt.Println("***第二轮结束后压测结果***")
		printTime(seEd-seSt, two)
		fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", three))
		<-time.After(10 * time.Second)
		wg.Add(int(three))
		for i := 0; i < int(three); i++ {
			//go TransactionHistory(third)
		}
		wg.Wait()
	}
}

func one1() {
	req := HttpRequest.NewRequest()
	data := make(map[string]interface{})
	data["end-time"] = strconv.FormatInt(time.Now().Unix(), 10)
	data["start_time"] = strconv.FormatInt(time.Now().AddDate(0, 0, -7).Unix(), 10)
	data["page"] = 1
	data["size"] = 10
	data["status"] = 0
	bs, _ := json.Marshal(&data)
	headers := make(map[string]string)
	headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", token, uid)
	res, _ := req.SetHeaders(headers).JSON().Post(orderUrl, bs)
	body, _ := res.Body()
	fmt.Println(string(body))
}
func OrderList(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT := time.Now().UnixNano() / 1e6
	for {
		req := HttpRequest.NewRequest()
		data := make(map[string]interface{})
		data["page"] = 1
		data["size"] = 10
		bs, _ := json.Marshal(&data)
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", token, uid)
		res, _ := req.SetHeaders(headers).JSON().Post(orderUrl, bs)
		defer res.Close()
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		lock.Lock()
		timeList = append(timeList, int(resTime))
		lock.Unlock()
		//channel <- resTime
		body, _ := res.Body()
		eT := time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var resMap map[string]interface{}
		err1 := json.Unmarshal(body, &resMap)
		if err1 != nil {
			log.Println("解析返回数据异常: ", err1)
		}
		if res.StatusCode() == 200 && resMap["status"] == 0 {
			atomic.AddInt64(&okNum, 1)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("接口返回：", string(body))
	}
}

func withDraws(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT := time.Now().UnixNano() / 1e6
	for {
		req := HttpRequest.NewRequest()
		data := make(map[string]interface{})
		data["end-time"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
		data["start_time"] = strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
		data["page"] = 1
		data["size"] = 10
		data["status"] = 0
		bs, _ := json.Marshal(&data)
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", token, uid)
		res, _ := req.SetHeaders(headers).JSON().Post(withDrawsUrl, bs)
		defer res.Close()
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		lock.Lock()
		timeList = append(timeList, int(resTime))
		lock.Unlock()
		//channel <- resTime
		body, _ := res.Body()
		eT := time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var resMap map[string]interface{}
		err1 := json.Unmarshal(body, &resMap)
		if err1 != nil {
			log.Println("解析返回数据异常: ", err1)
		}
		if res.StatusCode() == 200 && resMap["status"] == 0 {
			atomic.AddInt64(&okNum, 1)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("接口返回：", string(body))
	}
}
func getSpendList(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT := time.Now().UnixNano() / 1e6
	for {
		req := HttpRequest.NewRequest()
		data := make(map[string]interface{})
		data["etimestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
		data["stimestamp"] = strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
		data["page"] = 1
		data["size"] = 10
		data["coin_type"] = 0
		data["spend_type"] = 0
		bs, _ := json.Marshal(&data)
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", token, uid)
		res, _ := req.SetHeaders(headers).JSON().Post(getSpendListUrl, bs)
		defer res.Close()
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		lock.Lock()
		timeList = append(timeList, int(resTime))
		lock.Unlock()
		//channel <- resTime
		body, _ := res.Body()
		eT := time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var resMap map[string]interface{}
		err1 := json.Unmarshal(body, &resMap)
		if err1 != nil {
			log.Println("解析返回数据异常: ", err1)
		}
		if res.StatusCode() == 200 && resMap["status"] == 0 {
			atomic.AddInt64(&okNum, 1)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("接口返回：", string(body))
	}
}
