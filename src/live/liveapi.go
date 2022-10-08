package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
	"log"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	platForm  = 0
	loginType = 1
	//account         = "carl001"
	//pwd             = "carl001"
	host                  = "https://ycapi.mliveplus.com/"
	loginUrl              = "api/user/login"
	orderUrl              = "https://ycapi.mliveplus.com/api/charge/orderList"
	getActivity           = "http://ycapi.mliveplus.com/api/activity/getActivityById"
	getGiftUrl            = "http://ycapi.mliveplus.com/api/gift/getGiftList"
	getActivityListUrl    = "http://ycapi.mliveplus.com/api/activity/getActivityProfileList"
	getActivityRecordUrl  = "http://ycapi.mliveplus.com/api/activity/getActivityRecord"
	withDrawsUrl          = "https://ycapi.mliveplus.com/api/withdraw/withdraws"
	getSpendListUrl       = "https://ycapi.mliveplus.com/api/order/getSpendList"
	register1             = "api/user/regist"
	register2             = "api/user/register"
	atten                 = "http://ycapi.mliveplus.com/api/anchor/attentAnchor"
	userAddBank           = "api/user/addBank"
	userChangeMobile      = "api/user/changeMobile"
	userChangePayPwd      = "api/user/changePayPwd"
	userChangePwdUrl      = "api/user/changePwd"
	userCheckAnchorStatus = "api/user/checkAnchorStatus"
	userCheckSmsCode      = "api/user/checkSmsCode"
	userDefaultAvatar     = "api/user/defaultAvatar"
	userGetBanks          = "api/user/getBanks"
	userDelBanks          = "api/user/delBank"
	userGetFans           = "api/user/getFans"
	userGetInviteInfo     = "api/user/getInviteInfo"
	userGetUserAsset      = "api/user/getUserAsset"
	userGetUserInfo       = "api/user/getUserInfo"
	searchUrl             = "api/search"

	getExpertMenu            = "api/expert/getExpertCompMenu"
	getExpertFromMatch       = "api/expert/getExpertFromMatchList"
	getExpertPlanUrl         = "api/expert/getExpertPlan"
	getExpertPlanListUrl     = "api/expert/getExpertPlanList"
	favouriteMatchUrl        = "api/match/favoriteMatch"
	getCompMenuListURL       = "api/match/getCompMenuList"
	getHotMatchCardsUrl      = "api/match/getHotMatchCards4Crawler"
	getRecommendMatchCardUrl = "api/match/getRecommendMatches4Crawler"
	getSportMatchUrl         = "api/match/getSportMatch"
	getSportMatchCompUrl     = "api/match/getSportMatchComp"
	path                     = "/Users/eden/go/src/gostudy/src/ws/data.json"
)

var (
	//go:embed data.json
	f embed.FS
)

var (
	//创建计数器
	wg       = sync.WaitGroup{}
	sT       int64
	eT       int64
	useTime1 int64
	useTime2 int64
	useTime3 int64
	//num      int64 = 5 //设置并发数量
	okNum        int64 = 0 //初始化请求成功的数量
	firstOkNum   int64     //第一轮成功请求数
	secondOkNum  int64     //第二轮成功请求数
	thirdOkNum   int64     //第三轮成功请求数
	timeList     []int
	timeList1    []int //第一轮响应时间
	timeList2    []int //第二轮响应时间
	timeList3    []int //第三轮响应时间
	longTimeNum  int64 = 0
	longTimeNum1 int64
	longTimeNum2 int64
	longTimeNum3 int64
	channel            = make(chan int64)
	done               = make(chan struct{})
	lock               = sync.Mutex{}
	iNum         int64 = 0
	firstNum     int64 //第一轮请求数
	secondNum    int64 //第二轮请求数
	thirdNum     int64 //第三轮请求数
	one          int64 //第一轮并发数
	two          int64 //第二轮并发数
	three        int64 //第三轮并发数
	first        int64 //第一轮压测时长
	second       int64 //第二轮压测时长
	third        int64 //第三轮压测时长
	token        string
	uid          int
	phone        string
	userData     []map[string]string
	apiNum       int64

	hashChan  = make(chan string) //注册hash
	hashList  []string
	tokenChan = make(chan map[string]string) //登陆token
	transport *http.Transport
	bankId    = make(chan string)
	startTime string
)

type Response struct {
	Resp map[string]gjson.Result `json:"msg"`
}

//获取时间戳
func SetTs() string {
	return strconv.FormatInt(time.Now().Unix()*1000, 10)
}

func sumRespTime(timeList []int) int {
	sum := 0
	for _, index := range timeList {
		sum = sum + index
	}
	return sum
}

func maxRespTime(timeList []int) int {
	max := timeList[0]
	for _, index := range timeList {
		if index > max {
			max = index
		}
	}
	return max
}
func minRespTime(timeList []int) int {
	min := timeList[0]
	for _, index := range timeList {
		if index < min {
			min = index
		}
	}
	return min
}

func fiftyRespTime(timeList []int) int {
	sort.Ints(timeList)
	resSize := 0.5
	return timeList[int(float64(len(timeList))*resSize)-1]
}
func ninetyRespTime(timeList []int) int {
	sort.Ints(timeList)
	resSize := 0.9
	return timeList[int(float64(len(timeList))*resSize)-1]
}
func printTime(useTime int64, count int64, iNum int64, okNum int64, timeList []int, longTime int64) {
	fmt.Println("并发数: ", count)
	fmt.Println("请求数: ", iNum+count)
	fmt.Println("成功的数量：", len(timeList))
	fmt.Printf("\033[31m失败的数量：%v \033[0m \n", int(iNum+count)-len(timeList))
	fmt.Printf("\033[31m失败率：%.2f%v \033[0m \n", float64(int(iNum+count)-len(timeList))/float64(iNum+count)*100, "%")
	fmt.Println(fmt.Sprintf("耗时：%v秒", float64(useTime)/1000))
	fmt.Println("50%用户响应时间：" + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(fiftyRespTime(timeList)), float64(fiftyRespTime(timeList))/1e6))
	fmt.Println("90%用户响应时间：" + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(ninetyRespTime(timeList)), float64(ninetyRespTime(timeList))/1e6))
	fmt.Printf("\033[31m响应时间超过10秒的请求数：%v\033[0m \n", longTime)
	fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTime(timeList))/1e6)
	fmt.Printf("最小响应时间：%.3f微秒 ≈ %v秒 \n", float64(minRespTime(timeList)), float64(minRespTime(timeList))/1e6)
	fmt.Printf("平均响应时间是:%.3f秒 \n", float64(sumRespTime(timeList))/float64(len(timeList))/1e6)
	fmt.Printf("QPS：%.3f \n", float64(count)/(float64(sumRespTime(timeList))/float64(len(timeList))/1e6))
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
	fmt.Println(len(userData))
	apiList := []string{"(1)getExpertPlanList", "(2)getFans", "(3)getExpertCompMenu", "(4)register", "(5)getExpertFromMatchList", "(6)checkAnchorStatus", "(7)getExpertPlan", "(8)withdraws"}
	fmt.Println("接口列表：", apiList)
	fmt.Println("请选择并发接口序号(如：1): ")
	_, _ = fmt.Scan(&apiNum)
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
	startTimes := time.Now().UnixNano() / 1e6
	startTime = time.Now().Format("15：04：05")
	fmt.Printf("开始时间：%v \n", startTime)
	do()
	endTime := time.Now().UnixNano() / 1e6
	fmt.Printf("结束时间：%v \n", endTime)
	fmt.Printf("总耗时：%.3f 秒 \n", float64(endTime-startTimes)/1000-(10+10))
	//fmt.Println("总并发数:", one+two+three)
	//fmt.Println("总请求数: ", iNum)
	//fmt.Println("总成功的数量: ", firstOkNum+secondOkNum+thirdOkNum)
	//fmt.Printf("\033[31m总失败的数量: %v \033[0m \n", iNum-firstOkNum-secondOkNum-thirdOkNum)
	//fmt.Printf("\033[31m总失败率：%.2f%v \033[0m \n", float64(iNum-firstOkNum-secondOkNum-thirdOkNum)/float64(iNum)*100, "%")
	//fmt.Println("50%用户响应时间: " + fmt.Sprintf("%.3f秒", float64(fiftyRespTime())/1e6))
	//fmt.Println("90%用户响应时间: " + fmt.Sprintf("%.3f秒", float64(ninetyRespTime())/1e6))
	//fmt.Printf("最大响应时间: %.3f秒 \n", float64(maxRespTime())/1e6)
	//fmt.Printf("最小响应时间: %.3f微秒 约为%v秒 \n", float64(minRespTime()), float64(minRespTime())/1e6)
	//fmt.Printf("平均响应时间是:%.3f秒 \n", float64(sumRespTime())/float64(len(timeList))/1e6)
	//fmt.Printf("QPS: %.3f \n", float64(three)/(float64(sumRespTime())/float64(len(timeList))/1e6))
	runtime.GC()
	//_, _ = fmt.Scanf("h")
}

func do() {
	//wg.Add(1000)
	//for i := 0; i < 1000; i++ {
	//	go registerSetup()
	//}
	//wg.Wait()
	//time.Sleep(5 * time.Second)
	fmt.Println("第一轮压测开始...")
	fmt.Println(fmt.Sprintf("***首轮并发协程为%v***", one))
	wg.Add(int(one))
	for i := 0; i < int(one); i++ {
		switch apiNum {
		case 1:
			//go loginSetUp(i)
			go getUserInfo(first, i)
		case 2:
			//go loginSetUp(i)
			go getFans(first)
		case 3:
			//go loginSetUp(i)
			go getExpertCompMenu(first, i)
		case 4:
			go register(first)
		case 5:
			//go loginSetUp(i)
			go getExpertFromMatchList(first, i)
		case 6:
			//go loginSetUp(i)
			go UserCheckAnchorStatus(first, i)
		case 7:
			go getExpertPlan(first, i)
		case 8:
			withDraws(first)
		}
	}
	wg.Wait()
	firstNum = iNum
	firstOkNum = okNum
	timeList1 = timeList
	longTimeNum1 = longTimeNum
	useTime1 = eT - sT
	fmt.Println("第二轮压测开始...")
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", two))
	<-time.After(10 * time.Second)
	wg.Add(int(two))
	for i := 0; i < int(two); i++ {
		switch apiNum {
		case 1:
			//go loginSetUp(i)
			go getUserInfo(second, i)
		case 2:
			go loginSetUp(i)
			go UserChangeMobile(second, i)
		case 3:
			//go loginSetUp(i)
			go getExpertCompMenu(second, i)
		case 4:
			go register(second)
		case 5:
			//go loginSetUp(i)
			go getExpertFromMatchList(second, i)
		case 6:
			//go loginSetUp(i)
			go UserCheckAnchorStatus(second, i)
		case 7:
			go getExpertPlan(second, i)
		case 8:
			withDraws(second)
		}
	}
	wg.Wait()
	secondNum = iNum - firstNum
	secondOkNum = okNum - firstOkNum
	timeList2 = timeList[len(timeList1):]
	longTimeNum2 = longTimeNum - longTimeNum1
	useTime2 = eT - sT
	fmt.Println("第三轮压测开始...")
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", three))
	<-time.After(10 * time.Second)
	wg.Add(int(three))
	for i := 0; i < int(three); i++ {
		switch apiNum {
		case 1:
			//go loginSetUp(i)
			go getUserInfo(third, i)
		case 2:
			go loginSetUp(i)
			go UserChangeMobile(third, i)
		case 3:
			//go loginSetUp(i)
			go getExpertCompMenu(third, i)
		case 4:
			go register(third)
		case 5:
			//go loginSetUp(i)
			go getExpertFromMatchList(third, i)
		case 6:
			//go loginSetUp(i)
			go UserCheckAnchorStatus(third, i)
		case 7:
			go getExpertPlan(third, i)
		case 8:
			withDraws(third)
		}
	}
	wg.Wait()
	thirdNum = iNum - firstNum - secondNum
	thirdOkNum = okNum - firstOkNum - secondOkNum
	timeList3 = timeList[len(timeList1)+len(timeList2):]
	longTimeNum3 = longTimeNum - longTimeNum1 - longTimeNum2
	useTime3 = eT - sT
	fmt.Println("开始时间：", startTime)
	fmt.Println("\033[33m***第一轮压测结果***\033[0m")
	printTime(useTime1, one, firstNum, firstOkNum, timeList1, longTimeNum1)
	fmt.Println("\033[35m----------------------- \033[0m")
	fmt.Println("\033[33m***第二轮压测结果***\033[0m")
	printTime(useTime2, two, secondNum, secondOkNum, timeList2, longTimeNum2)
	fmt.Println("\033[35m----------------------- \033[0m")
	fmt.Println("\033[33m***第三轮压测结果***\033[0m")
	printTime(useTime3, three, thirdNum, thirdOkNum, timeList3, longTimeNum3)
	fmt.Println("\033[35m----------------------- \033[0m")

}

//添加银行卡
func UserAddBank(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		fmt.Println("获取的chanMap: ", chMap)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		req := HttpRequest.NewRequest()
		rand.Seed(time.Now().UnixNano())
		card := rand.Intn((1000 - 1) + 1)
		card2 := rand.Intn((1000 - 1) + 1)
		payLoad := make(map[string]interface{})
		payLoad["bank"] = "工商银行"
		payLoad["card"] = fmt.Sprintf("%v%v", card, card2)
		payLoad["id_card"] = fmt.Sprintf("%v%v", card, card2)
		payLoad["name"] = "eden"
		payLoad["smscode"] = "999999"
		res, err2 := req.SetHeaders(headers).JSON().Post(host+userAddBank, payLoad)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("请求Url: ", userAddBank)
		body, _ := res.Body()
		log.Println("接口返回：", string(body))
		res.Close()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常：", string(body))
		}
	}
}

//修改手机号
func UserChangeMobile(times int64, i int) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	//chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		req := HttpRequest.NewRequest()
		rand.Seed(time.Now().UnixNano())
		//m1 := rand.Intn((10000 - 1000) + 1000)
		//m2 := rand.Intn((10000 - 1000) + 1000)
		payLoad := make(map[string]interface{})
		payLoad["new_mobile"] = userData[i]["phone"]
		payLoad["new_smscode"] = "999999"
		payLoad["ori_sms_code"] = "999999"
		res, err2 := req.SetHeaders(headers).JSON().Post(host+userChangeMobile, payLoad)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("请求Url: ", userChangeMobile)
		body, _ := res.Body()
		log.Println("接口返回：", string(body))
		res.Close()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常：", string(body))
		}
	}
}

//修改支付密码
func UserChangePayPwd(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		req := HttpRequest.NewRequest()
		rand.Seed(time.Now().UnixNano())
		//m1 := rand.Intn((10000 - 1000) + 1000)
		//m2 := rand.Intn((10000 - 1000) + 1000)
		payLoad := make(map[string]interface{})
		payLoad["pwd"] = "111111"
		payLoad["smscode"] = "999999"
		res, err2 := req.SetHeaders(headers).JSON().Post(host+userChangePayPwd, payLoad)
		if err2 != nil {
			log.Println("\033[31m请求异常：\033[0m", err2)
			continue
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("请求Url: ", userChangePayPwd)
		body, _ := res.Body()
		log.Println("接口返回：", string(body))
		res.Close()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		if int(resTime) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常：", string(body))
		}
	}
}

//修改密码
func UserChangePwd(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		req := HttpRequest.NewRequest()
		rand.Seed(time.Now().UnixNano())
		//m1 := rand.Intn((10000 - 1000) + 1000)
		//m2 := rand.Intn((10000 - 1000) + 1000)
		payLoad := make(map[string]interface{})
		payLoad["pwd"] = chMap["pwd"]
		payLoad["smscode"] = "999999"
		res, err2 := req.SetHeaders(headers).JSON().Post(host+userChangePwdUrl, payLoad)
		if err2 != nil {
			log.Println("\033[31m请求异常：\033[0m", err2)
			continue
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("请求Url: ", userChangePwdUrl)
		body, _ := res.Body()
		log.Println("接口返回：", string(body))
		res.Close()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		if int(resTime) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("\033[31m响应异常：\033[0m", string(body))
		}
	}
}

//检查主播状态
func UserCheckAnchorStatus(times int64, i int) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	//chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		req := HttpRequest.NewRequest()
		res, err2 := req.SetHeaders(headers).JSON().Post(host + userCheckAnchorStatus)
		if err2 != nil {
			log.Println("\033[31m请求异常：\033[0m", err2)
			continue
		}
		//log.Println("响应码：", res.StatusCode())
		//log.Println("请求Url: ", host+userCheckAnchorStatus)
		body, _ := res.Body()
		//log.Println("接口返回：", string(body))
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		if int(resTime) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		log.Printf("响应时间：%v毫秒", resTime)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("\033[31m响应异常：\033[0m", string(body))
		}
		res.Close()
	}
}

//校验验证码
func UserCheckSmsCode(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		req := HttpRequest.NewRequest()
		rand.Seed(time.Now().UnixNano())
		//m1 := rand.Intn((10000 - 1000) + 1000)
		//m2 := rand.Intn((10000 - 1000) + 1000)
		payLoad := make(map[string]interface{})
		payLoad["mobile"] = chMap["phone"]
		payLoad["smscode"] = "999999"
		res, err2 := req.SetHeaders(headers).JSON().Post(host+userCheckSmsCode, payLoad)
		if err2 != nil {
			log.Println("\033[31m请求异常：\033[0m", err2)
			continue
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("请求Url: ", userCheckSmsCode)
		body, _ := res.Body()
		log.Println("接口返回：", string(body))
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		if int(resTime) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("\033[31m响应异常：\033[0m", string(body))
		}
		res.Close()
	}
}

//用户默认头像
func UserDefaultAvatar(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		req := HttpRequest.NewRequest()
		rand.Seed(time.Now().UnixNano())
		//m1 := rand.Intn((10000 - 1000) + 1000)
		//m2 := rand.Intn((10000 - 1000) + 1000)
		res, err2 := req.SetHeaders(headers).JSON().Post(host + userDefaultAvatar)
		if err2 != nil {
			log.Println("\033[31m请求异常：\033[0m", err2)
			continue
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("请求Url: ", userDefaultAvatar)
		body, _ := res.Body()
		log.Println("接口返回：", string(body))
		res.Close()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		if int(resTime) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("\033[31m响应异常：\033[0m", string(body))
		}
	}
}

func getFans(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	//chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", "hj+TzZwPBcltZ4mDOLI9bdeDkr6WtE0y", 6023732)
		req := HttpRequest.NewRequest()
		res, err2 := req.SetHeaders(headers).JSON().Post(host + userGetFans)
		if err2 != nil {
			log.Println("\033[31m请求异常：\033[0m", err2)
			continue
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("请求Url: ", host+userGetFans)
		body, _ := res.Body()
		log.Println("接口返回：", string(body))
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		if int(resTime) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("\033[31m响应异常：\033[0m", string(body))
		}
		res.Close()
	}
}

func getInviteInfo(times int64, i int) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	//chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		req := HttpRequest.NewRequest()
		rand.Seed(time.Now().UnixNano())
		//m1 := rand.Intn((10000 - 1000) + 1000)
		//m2 := rand.Intn((10000 - 1000) + 1000)
		res, err2 := req.SetHeaders(headers).JSON().Post(host + userGetInviteInfo)
		if err2 != nil {
			log.Println("\033[31m请求异常：\033[0m", err2)
		}
		//log.Println("响应码：", res.StatusCode())
		//log.Println("请求Url: ", userGetInviteInfo)
		body, _ := res.Body()
		//log.Println("接口返回：", string(body))
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		if int(resTime) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		log.Printf("响应时间：%v毫秒", resTime)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("\033[31m响应异常：\033[0m", string(body))
		}
		res.Close()
	}
}

func getUserAsset(times int64, i int) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	//chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		req := HttpRequest.NewRequest()
		res, err2 := req.SetHeaders(headers).JSON().Post(host + userGetUserAsset)
		if err2 != nil {
			log.Println("\033[31m请求异常：\033[0m", err2)
			continue
		}
		//log.Println("响应码：", res.StatusCode())
		//log.Println("请求Url: ", userGetUserAsset)
		body, _ := res.Body()
		//log.Println("接口返回：", string(body))
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		if int(resTime) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		log.Printf("响应时间：%v毫秒", resTime)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("\033[31m响应异常：\033[0m", string(body))
		}
		res.Close()
	}
}

func getUserInfo(times int64, i int) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	//chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		req := HttpRequest.NewRequest()
		res, err2 := req.SetHeaders(headers).JSON().Post(host + userGetUserInfo)
		if err2 != nil {
			log.Println("\033[31m请求异常：\033[0m", err2)
		}
		//log.Println("响应码：", res.StatusCode())
		//log.Println("请求Url: ", userGetUserInfo)
		body, _ := res.Body()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		log.Println("接口返回：", string(body))
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		if int(resTime) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		log.Printf("响应时间：%v毫秒", resTime)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("\033[31m响应异常：\033[0m", string(body))
		}
		res.Close()
	}
}

func search(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		data := make(map[string]interface{})
		data["keyword"] = "美乐"
		data["type"] = 0
		id, _ := strconv.Atoi(userData[i]["id"])
		data["user_id"] = id
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host+searchUrl, data)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		log.Println("响应码：", res.StatusCode())
		log.Println("接口返回：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && resMap["msg"] == "success" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		log.Println("请求Url: ", host+searchUrl)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}

}

func getExpertCompMenu(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host + getExpertMenu)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		log.Println("响应码：", res.StatusCode())
		//log.Println("接口返回：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && string(body) != "" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		log.Println("请求Url: ", host+getExpertMenu)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}

}

func getExpertFromMatchList(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		data := make(map[string]interface{})
		data["match_id"] = 3672203
		data["page"] = 1
		data["size"] = 10
		data["sport_id"] = 2
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host+getExpertFromMatch, data)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		log.Println("响应码：", res.StatusCode())
		log.Println("接口返回：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && string(body) != "" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		log.Println("请求Url: ", host+getExpertFromMatch)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}

}

func getExpertPlan(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		data := make(map[string]interface{})
		data["id"] = 46
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host+getExpertPlanUrl, data)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		//log.Println("响应码：", res.StatusCode())
		//log.Println("接口返回：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && string(body) != "" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		//log.Println("请求Url: ", host+getExpertPlanUrl)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}

}

func getExpertPlanList(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		data := make(map[string]interface{})
		data["expert_id"] = 0
		freeTypes := []int{0}
		data["fee_types"] = freeTypes
		data["match_id"] = 3672203
		data["page"] = 1
		data["size"] = 10
		data["info_expert_id"] = 0
		playTypes := []int{0}
		data["play_types"] = playTypes
		data["sport_id"] = 2
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host+getExpertPlanListUrl, data)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		//log.Println("响应码：", res.StatusCode())
		//log.Println("接口返回：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && string(body) != "" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		//log.Println("请求Url: ", host+getExpertPlanListUrl)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}
}

func favoriteMatch(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	var z = 0
	for {
		z++
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		data := make(map[string]interface{})
		data["match_id"] = 3779692
		//data["user_id"] = userData[i]["id"]
		data["sport_id"] = 1
		if z%2 == 0 {
			data["operate"] = 1
		} else {
			data["operate"] = 2
		}
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host+favouriteMatchUrl, data)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		//log.Println("响应：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && string(body) != "" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		//log.Println("请求Url: ", host+getExpertPlanListUrl)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}
}

func getCompMenuList(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		data := make(map[string]interface{})
		data["page"] = 1
		data["size"] = 10
		data["date"] = time.Now().Format("2006-01-02")
		data["match_status"] = 0
		data["sport_id"] = 1
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host+getCompMenuListURL, data)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		//log.Println("响应：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && string(body) != "" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		//log.Println("请求Url: ", host+getExpertPlanListUrl)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}
}

func getHotMatchCards(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host + getHotMatchCardsUrl)
		if err2 != nil {
			log.Println("请求异常：", err2)
			continue
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		//log.Println("响应：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && string(body) != "" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		//log.Println("请求Url: ", host+getExpertPlanListUrl)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}
}

func getRecommendMatchCard(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		data := make(map[string]interface{})
		data["sport_id"] = 1
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host+getRecommendMatchCardUrl, data)
		if err2 != nil {
			log.Println("请求异常：", err2)
			continue
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		//log.Println("响应：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && string(body) != "" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		//log.Println("请求Url: ", host+getExpertPlanListUrl)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}
}

func getSportMatch(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		data := make(map[string]interface{})
		//data["comp"] = "热门"
		data["sport_id"] = 1
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host+getSportMatchUrl, data)
		if err2 != nil {
			log.Println("请求异常：", err2)
			continue
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		//log.Println("响应：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && string(body) != "" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		//log.Println("请求Url: ", host+getExpertPlanListUrl)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}
}

func getSportMatchComp(times int64, i int) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		//chMap := <-ch
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", userData[i]["token"], userData[i]["id"])
		data := make(map[string]interface{})
		data["comp"] = "热门"
		data["sport_id"] = 1
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(host+getSportMatchCompUrl, data)
		if err2 != nil {
			log.Println("请求异常：", err2)
			continue
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		//log.Println("响应：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && string(body) != "" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		//log.Println("请求Url: ", host+getExpertPlanListUrl)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		if int(e-s) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}
}

//银行卡列表
func UserGetBanks(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		req := HttpRequest.NewRequest()
		rand.Seed(time.Now().UnixNano())
		//m1 := rand.Intn((10000 - 1000) + 1000)
		//m2 := rand.Intn((10000 - 1000) + 1000)
		res, err2 := req.SetHeaders(headers).JSON().Post(host + userGetBanks)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("请求Url: ", userGetBanks)
		body, _ := res.Body()
		log.Println("接口返回：", string(body))
		res.Close()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)

		//channelLive <- resTime * 1000
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常：", string(body))
		}
	}
}

func UserGetBanksSetUp(ch chan map[string]string) {
	chMap := <-ch
	headers := make(map[string]string)
	headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
	req := HttpRequest.NewRequest()
	rand.Seed(time.Now().UnixNano())
	//m1 := rand.Intn((10000 - 1000) + 1000)
	//m2 := rand.Intn((10000 - 1000) + 1000)
	res, err2 := req.SetHeaders(headers).JSON().Post(host + userGetBanks)
	if err2 != nil {
		log.Println("\033[31m请求异常：\033[0m", err2)
	}
	body, _ := res.Body()
	res.Close()
	bankId <- gjson.ParseBytes(body).Map()["data"].Array()[0].Map()["id"].String()
}

func UserDelBanks(times int64, ch chan map[string]string, chBankId chan string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		req := HttpRequest.NewRequest()
		rand.Seed(time.Now().UnixNano())
		//m1 := rand.Intn((10000 - 1000) + 1000)
		//m2 := rand.Intn((10000 - 1000) + 1000)
		res, err2 := req.SetHeaders(headers).JSON().Post(host + userDelBanks + "/" + <-chBankId)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("请求Url: ", userDelBanks)
		body, _ := res.Body()
		log.Println("接口返回：", string(body))
		res.Close()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000

		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常：", string(body))
		}
	}
}

func registerSetup() {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	req := HttpRequest.NewRequest()
	payLoad := make(map[string]string)
	rand.Seed(time.Now().UnixNano())
	payLoad["mobile"] = fmt.Sprintf("18%v6%v", rand.Intn(9999-1000)+1000, rand.Intn(9999-1000)+1000)
	payLoad["smscode"] = "999999"
	res, err2 := req.JSON().Post(host+register1, payLoad)
	if err2 != nil {
		log.Println("\033[31m获取注册hash异常：\033[0m", err2)
	}
	body, _ := res.Body()
	log.Println("获取注册hash：", string(body))
	var resMap map[string]interface{}
	err := json.Unmarshal(body, &resMap)
	if err != nil {
		fmt.Println(err)
	}
	hash := gjson.ParseBytes(body).Map()["data"].Map()["mobile_hash"].String()
	lock.Lock()
	hashList = append(hashList, hash)
	lock.Unlock()
	phone = payLoad["mobile"]
	res.Close()

}
func register(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	var a = 0
	for {
		a++
		req := HttpRequest.NewRequest()
		rand.Seed(time.Now().UnixNano())
		payLoad := make(map[string]interface{})
		ranS := rand.Intn(9999-1000) + 1000
		ranS1 := rand.Intn(9999-1000) + 1000
		payLoad["account"] = fmt.Sprintf("xy%v3%v", ranS, ranS1)
		payLoad["mobile_hash"] = hashList[a]
		payLoad["invite_code"] = ""
		payLoad["platform"] = 0
		payLoad["sex"] = 0
		payLoad["pwd"] = fmt.Sprintf("xy%v5%v6", ranS, ranS1)
		res, err2 := req.JSON().Post(host+register2, payLoad)
		if err2 != nil {
			log.Println("\033[31m注册异常：\033[0m", err2)
			continue
		}
		body, _ := res.Body()
		log.Println("注册返回：", string(body))
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		if int(resTime) > 10000 {
			atomic.AddInt64(&longTimeNum, 1)
		}
		log.Printf("响应时间：%v毫秒", resTime)
		//channelLive <- resTime * 1000
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("\033[31m响应异常：\033[0m", string(body))
		}
		res.Close()
	}

	//fmt.Println(string(body))
	//account = gjson.ParseBytes(body).Map()["data"].Map()["account"].String()
	//pwd = gjson.ParseBytes(body).Map()["data"].Map()["account"].String()
}

func attenAnchor(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		req := HttpRequest.NewRequest()
		payLoad := make(map[string]interface{})
		payLoad["anchorid"] = 27236
		res, err2 := req.SetHeaders(headers).JSON().Post(atten, payLoad)
		if err2 != nil {
			log.Println("关注异常：", err2)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("请求Url: ", atten)
		body, _ := res.Body()
		log.Println("接口返回：", string(body))
		res.Close()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		//channelLive <- resTime * 1000
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		var resMap map[string]interface{}
		err := json.Unmarshal(body, &resMap)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常：", string(body))
		}
	}

	//fmt.Println(string(body))
	//account = gjson.ParseBytes(body).Map()["data"].Map()["account"].String()
	//pwd = gjson.ParseBytes(body).Map()["data"].Map()["account"].String()
}
func login(times int64, i int) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		req := HttpRequest.NewRequest()
		payLoad := make(map[string]interface{})
		payLoad["account"] = userData[i]["account"]
		payLoad["login_type"] = loginType
		payLoad["platform"] = platForm
		payLoad["pwd"] = userData[i]["pwd"]
		payLoad["source_id"] = 0
		res, err2 := req.JSON().Post(host+loginUrl, payLoad)
		if err2 != nil {
			log.Println("请求异常：", err2)
			continue
		}
		body, _ := res.Body()
		//log.Println("登陆返回：", string(body))
		res.Close()
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", resTime)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
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
	}

	//token = gjson.ParseBytes(body).Map()["data"].Map()["token"].String()
	//uid = int(gjson.ParseBytes(body).Map()["data"].Map()["id"].Int())
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
	//payLoad["account"] = "den15968"
	payLoad["login_type"] = loginType
	payLoad["platform"] = platForm
	payLoad["pwd"] = userData[i]["pwd"]
	payLoad["source_id"] = 0
	res, _ := req.JSON().Post(reqUrl, payLoad)
	defer res.Close()
	body, _ := res.Body()
	fmt.Println("登陆：", string(body))
	chMap := make(map[string]string)
	chMap["token"] = gjson.ParseBytes(body).Map()["data"].Map()["token"].String()
	chMap["id"] = gjson.ParseBytes(body).Map()["data"].Map()["id"].String()
	chMap["nick_name"] = gjson.ParseBytes(body).Map()["data"].Map()["nick_name"].String()
	chMap["avatar"] = gjson.ParseBytes(body).Map()["data"].Map()["avatar"].String()
	chMap["pwd"] = gjson.ParseBytes(body).Map()["data"].Map()["account"].String()
	//chMap["phone"] = userData[i]["phone"]
	tokenChan <- chMap
	fmt.Println("chanMap: ", chMap)
}

func getActivityById(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		req := HttpRequest.NewRequest()
		data := make(map[string]interface{})
		data["id"] = 17
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		res, err2 := req.SetHeaders(headers).JSON().Post(getActivity, data)
		log.Println("响应码：", res.StatusCode())
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", resTime)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		//channel <- resTime * 1000
		body, err3 := res.Body()
		if err3 != nil {
			log.Println("解析响应body异常：", err3)
		}
		log.Println("接口返回：", string(body))
		eT = time.Now().UnixNano() / 1e6
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
		res.Close()
	}
}

func getActivityList(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		req := HttpRequest.NewRequest()
		data := make(map[string]interface{})
		data["page"] = 1
		data["size"] = 10
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		res, err2 := req.SetHeaders(headers).JSON().Post(getActivityListUrl, data)
		log.Println("响应码：", res.StatusCode())
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", resTime)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		//channel <- resTime * 1000
		body, err3 := res.Body()
		if err3 != nil {
			log.Println("解析响应body异常：", err3)
		}
		log.Println("接口返回：", string(body))
		eT = time.Now().UnixNano() / 1e6
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
		res.Close()
	}
}

func getActivityRecord(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	chMap := <-ch
	//sT = time.Now().UnixNano() / 1e6
	req := HttpRequest.NewRequest()
	headers := make(map[string]string)
	headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
	res, err2 := req.SetHeaders(headers).JSON().Post(getActivityRecordUrl)
	log.Println("响应码：", res.StatusCode())
	if err2 != nil {
		log.Println("请求异常：", err2)
	}
	resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
	log.Printf("响应时间：%v毫秒", resTime)
	lock.Lock()
	timeList = append(timeList, int(resTime*1000))
	lock.Unlock()
	//channel <- resTime * 1000
	body, err3 := res.Body()
	if err3 != nil {
		log.Println("解析响应body异常：", err3)
	}
	log.Println("接口返回：", string(body))
	//eT = time.Now().UnixNano() / 1e6
	atomic.AddInt64(&iNum, 1)
	var resMap map[string]interface{}
	err1 := json.Unmarshal(body, &resMap)
	if err1 != nil {
		log.Println("解析返回数据异常: ", err1)
	}
	if res.StatusCode() == 200 && resMap["status"] == 0 {
		atomic.AddInt64(&okNum, 1)
	}
	defer res.Close()

}

func getGift(times int64) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	for {
		//chMap := <-ch
		sT = time.Now().UnixNano() / 1e6
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", "hj+TzZwPBcnWozFbMrTqg+HGeiAYYdAG", 6023717)
		s := time.Now().UnixNano() / 1e6
		res, err2 := req.SetHeaders(headers).JSON().Post(getGiftUrl)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		e := time.Now().UnixNano() / 1e6
		body, _ := res.Body()
		log.Println("响应码：", res.StatusCode())
		log.Println("接口返回：", string(body))
		var resMap map[string]interface{}
		_ = json.Unmarshal(body, &resMap)
		if res.StatusCode() == 200 && resMap["status"] == 0 && resMap["msg"] == "success" {
			atomic.AddInt64(&okNum, 1)
		}
		res.Close()
		log.Println("请求Url: ", getGiftUrl)
		//resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", e-s)
		lock.Lock()
		timeList = append(timeList, int((e-s)*1000))
		lock.Unlock()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
	}

}
func orderList(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		req := HttpRequest.NewRequest()
		data := make(map[string]interface{})
		data["page"] = 1
		data["size"] = 10
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		res, _ := req.SetHeaders(headers).JSON().Post(orderUrl, data)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		//channel <- resTime * 1000
		body, _ := res.Body()
		res.Close()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var resMap map[string]interface{}
		err1 := json.Unmarshal(body, &resMap)
		if err1 != nil {
			log.Println("解析返回数据异常: ", err1)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("接口返回：", string(body))
	}
}
func init() {
	transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 40 * time.Second,
			//KeepAlive: 60 * time.Second,
		}).DialContext,
		//MaxIdleConns: 1000000,
		//IdleConnTimeout:       90 * time.Second,
		//TLSHandshakeTimeout:   5 * time.Second,
		//ExpectContinueTimeout: 1 * time.Second,
	}
}

func withDraws(times int64) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	//chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		req := HttpRequest.NewRequest()
		data := make(map[string]interface{})
		data["end-time"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
		data["start_time"] = strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
		data["page"] = 1
		data["size"] = 10
		data["status"] = 0
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", "hj+TzZwPBckc56QGx+QjDiKdntUEtIkl", "6028700")
		res, _ := req.SetHeaders(headers).JSON().Post(withDrawsUrl, data)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		//channel <- resTime
		body, _ := res.Body()
		res.Close()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var resMap map[string]interface{}
		err1 := json.Unmarshal(body, &resMap)
		if err1 != nil {
			log.Println("解析返回数据异常: ", err1)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("接口返回：", string(body))
	}
}
func getSpendList(times int64, ch chan map[string]string) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	chMap := <-ch
	sT = time.Now().UnixNano() / 1e6
	for {
		req := HttpRequest.NewRequest()
		data := make(map[string]interface{})
		data["etimestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
		data["stimestamp"] = strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
		data["page"] = 1
		data["size"] = 10
		data["coin_type"] = 0
		data["spend_type"] = 0
		headers := make(map[string]string)
		headers["authorization"] = fmt.Sprintf("token=%v; uid=%v", chMap["token"], chMap["id"])
		res, _ := req.SetHeaders(headers).JSON().Post(getSpendListUrl, data)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		//channel <- resTime
		body, _ := res.Body()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		res.Close()
		atomic.AddInt64(&iNum, 1)
		var resMap map[string]interface{}
		err1 := json.Unmarshal(body, &resMap)
		if err1 != nil {
			log.Println("解析返回数据异常: ", err1)
		}
		if res.StatusCode() == 200 && int(resMap["status"].(float64)) == 0 {
			atomic.AddInt64(&okNum, 1)
		}
		log.Println("响应码：", res.StatusCode())
		log.Println("接口返回：", string(body))
	}
}
