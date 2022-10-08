package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	wsLiveUri    = "ws://chat.mliveplus.com/"
	wsType       = "login"
	sendDataType = "say"
	sendContent  = "斗皇强者，恐怖如斯！😂😂😂😂"
	content2     = "仙路尽头谁为峰，一遇无始道成空!🤙🤙🤙🤙🤙"
	content3     = "天不生我李淳刚，剑道万古如长夜!😭😭😭😭"
	content4     = "吾儿王腾有大帝之资～！😡😡😡😡"
	content5     = "遮天，谁以弹指而遮天？天地不仁，何以万物争相残？蚍蜉撼树可献命，是因其心蛇吞仙！🌚🌚🌚🌚🌚🌚"
	content6     = "黑袍老者目中精光一闪，旋即倒吸一口凉气！👽👽👽"
	roomId       = 2
	platForm     = 0
	loginType    = 1
	loginUrl     = "https://api.mliveplus.com/api/user/login"
	getLiveUrl   = "https://api.mliveplus.com/webapi/live/getLivePageData"
	path         = "/Users/eden/go/src/gostudy/src/ws/data.json"
)

var (
	//go:embed data.json
	f embed.FS
)

var (
	wgLive = sync.WaitGroup{}
	//num         int64 = 2
	sTime            int64
	eTime            int64
	useTime1         int64
	useTime2         int64
	useTime3         int64
	okNumLive        int64 = 0
	firstOkNumLive   int64 //第一轮请求成功数
	secondOkNumLive  int64 //第二轮请求成功数
	thirdOkNumLive   int64 //第三轮请求成功数
	channelLive      = make(chan int64)
	resTimeLiveList  []int
	resTimeLiveList1 []int //第一轮响应时间
	resTimeLiveList2 []int //第二轮响应时间
	resTimeLiveList3 []int //第三轮响应时间
	lockLive               = sync.Mutex{}
	iLive            int64 = 0
	iLiveFirst       int64 //第一轮请求数
	iLiveSecond      int64 //第二轮请求数
	iLiveThird       int64 //第三轮请求数
	doneLive1        = make(chan struct{})
	//resTime     int64
	oneLive        int64 //第一轮并发数
	twoLive        int64 //第二轮并发数
	threeLive      int64 //第三轮并发数
	firstLive      int64 //第一轮运行时长
	secondLive     int64 //第二轮运行时长
	thirdLive      int64 //第三轮运行时长
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

	hashChan  = make(chan string)                 //注册hash
	tokenChan = make(chan map[string]interface{}) //登陆token
)

func maxRespTimeLive(resTimeLiveList []int) int {
	max := resTimeLiveList[0]
	for _, i := range resTimeLiveList {
		if i > max {
			max = i
		}
	}
	return max
}
func minRespTimeLive(resTimeLiveList []int) int {
	min := resTimeLiveList[0]
	for _, i := range resTimeLiveList {
		if i < min {
			min = i
		}
	}
	return min
}
func sumResTimeLive(resTimeLiveList []int) int {
	sum := 0
	for _, i := range resTimeLiveList {
		sum += i
	}
	return sum
}
func fiftyRespTime(resTimeLiveList []int) int {
	sort.Ints(resTimeLiveList)
	resSize := 0.5
	return resTimeLiveList[int(float64(len(resTimeLiveList))*resSize)-1]
}
func ninetyRespTime(resTimeLiveList []int) int {
	sort.Ints(resTimeLiveList)
	resSize := 0.9
	return resTimeLiveList[int(float64(len(resTimeLiveList))*resSize)-1]
}

func printTimeLive(usetime int64, count int64, iLive int64, okNumLive int64, resTimeLiveList []int) {
	fmt.Println(len(resTimeLiveList))
	fmt.Println("并发数：", count)
	fmt.Println("请求数: ", iLive)
	fmt.Println("成功的数量：", len(resTimeLiveList))
	fmt.Printf("\033[31m失败的数量：%v\033[0m \n", int(iLive)-len(resTimeLiveList))
	fmt.Printf("\033[31m失败率：%.2f%v\033[0m \n", float64(int(iLive)-len(resTimeLiveList))/float64(iLive)*100, "%")
	fmt.Println(fmt.Sprintf("耗时: %v秒", float64(usetime)/1000))
	fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTimeLive(resTimeLiveList))/1e6)
	fmt.Printf("最小响应时间：%.3f微秒 ≈ %v秒 \n", float64(minRespTimeLive(resTimeLiveList)), float64(minRespTimeLive(resTimeLiveList))/1e6)
	fmt.Println("50%用户响应时间: " + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(fiftyRespTime(resTimeLiveList)), float64(fiftyRespTime(resTimeLiveList))/1e6))
	fmt.Println("90%用户响应时间: " + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(ninetyRespTime(resTimeLiveList)), float64(ninetyRespTime(resTimeLiveList))/1e6))
	fmt.Printf("平均响应时间是:%.3f秒 \n", float64(sumResTimeLive(resTimeLiveList))/float64(iLive)/1e6)
	fmt.Printf("QPS：%.3f \n", float64(count)/(float64(sumResTimeLive(resTimeLiveList))/float64(iLive)/1e6))
}

func runLive() {
	fmt.Println("第一轮压测开始...")
	fmt.Println(fmt.Sprintf("***首轮并发用户为%v协程***", oneLive))
	wgLive.Add(int(oneLive))
	for i := 0; i < int(oneLive); i++ {
		//go loginSetUp(i)
		go liveSend(firstLive, i)
	}
	wgLive.Wait()
	iLiveFirst = iLive
	firstOkNumLive = okNumLive
	resTimeLiveList1 = resTimeLiveList
	useTime1 = eTime - sTime
	fmt.Println("第二轮压测开始...")
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", twoLive))
	<-time.After(10 * time.Second)
	wgLive.Add(int(twoLive))
	for i := 0; i < int(twoLive); i++ {
		//go loginSetUp(i)
		go liveSend(secondLive, i)
	}
	wgLive.Wait()
	iLiveSecond = iLive - iLiveFirst
	secondOkNumLive = okNumLive - firstOkNumLive
	resTimeLiveList2 = resTimeLiveList[len(resTimeLiveList1):]
	useTime2 = eTime - sTime
	fmt.Println("第三轮压测开始...")
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", threeLive))
	<-time.After(10 * time.Second)
	wgLive.Add(int(threeLive))
	for i := 0; i < int(threeLive); i++ {
		//go loginSetUp(i)
		go liveSend(thirdLive, i)
	}
	wgLive.Wait()
	iLiveThird = iLive - iLiveFirst - iLiveSecond
	thirdOkNumLive = okNumLive - firstOkNumLive - secondOkNumLive
	resTimeLiveList3 = resTimeLiveList[len(resTimeLiveList1)+len(resTimeLiveList2):]
	useTime3 = eTime - sTime
	fmt.Println("\033[33m***第一轮压测结果***\033[0m")
	printTimeLive(useTime1, oneLive, iLiveFirst, firstOkNumLive, resTimeLiveList1)
	fmt.Println("\033[35m----------------------- \033[0m")
	fmt.Println("\033[33m***第二轮压测结果***\033[0m")
	printTimeLive(useTime2, twoLive, iLiveSecond, secondOkNumLive, resTimeLiveList2)
	fmt.Println("\033[35m----------------------- \033[0m")
	fmt.Println("\033[33m***第三轮压测结果***\033[0m")
	printTimeLive(useTime3, threeLive, iLiveThird, thirdOkNumLive, resTimeLiveList3)
	fmt.Println("\033[35m----------------------- \033[0m")
}

func main() {

	//login()
	//byteData1, err := ioutil.ReadFile(path)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err2 := json.Unmarshal(byteData1, &userData)
	//if err1 != nil {
	//	fmt.Println(err2)
	//}

	//for i := 0; i < 1; i++ {
	//	a1()
	//}
	//byteData, err := ioutil.ReadFile(path)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//var s []map[string]string
	//err1 := json.Unmarshal(byteData, &s)
	//if err1 != nil {
	//	fmt.Println(err1)
	//}
	//s = append(s, accountList...)
	//fmt.Println(len(s))
	//for index, i := range s {
	//	if i["account"] == "" {
	//		s = append(s[:index], s[index+1:]...)
	//	}
	//}
	//f, _ := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0666)
	//bS, _ := json.Marshal(s)
	//_, _ = f.Write(bS)

	//flag.Int64Var(&oneLive, "oneLiveCount", 10000, "第一轮并发数")
	//flag.Int64Var(&firstLive, "oneLiveTime", 100, "第一轮压测运行时长(秒)")
	//flag.Int64Var(&twoLive, "twoLiveCount", 15000, "第二轮并发数")
	//flag.Int64Var(&secondLive, "twoLiveTime", 100, "第二轮压测运行时长(秒)")
	//flag.Int64Var(&threeLive, "threeLiveCount", 20000, "第三轮并发数")
	//flag.Int64Var(&thirdLive, "threeLiveTime", 100, "第三轮压测运行时长(秒)")
	//flag.Parse()

	//byteData1, err := f.ReadFile("data.json")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err1 := json.Unmarshal(byteData1, &userData)
	//if err1 != nil {
	//	fmt.Println(err1)
	//}

	byteData, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	err1 := json.Unmarshal(byteData, &userData)
	if err1 != nil {
		fmt.Println(err1)
	}

	fmt.Println("第一轮压测用户数：")
	_, _ = fmt.Scan(&oneLive)
	fmt.Println("第一轮运行时长(秒)：")
	_, _ = fmt.Scan(&firstLive)
	fmt.Println("第二轮压测用户数：")
	_, _ = fmt.Scan(&twoLive)
	fmt.Println("第二轮运行时长(秒): ")
	_, _ = fmt.Scan(&secondLive)
	fmt.Println("第三轮压测用户数：")
	_, _ = fmt.Scan(&threeLive)
	fmt.Println("第三轮运行时长(秒): ")
	_, _ = fmt.Scan(&thirdLive)
	s1 := time.Now().UnixNano() / 1e6
	runLive()
	e1 := time.Now().UnixNano() / 1e6
	fmt.Println("**********")
	fmt.Println(fmt.Sprintf("总耗时：%v秒", float64(e1-s1)/1000-(10+10)))
	//fmt.Println("总并发数：", oneLive+twoLive+threeLive)
	//fmt.Println("总请求数: ", iLive-(oneLive+twoLive+threeLive))
	//fmt.Printf("成功的数量：%v \n", okNumLive)
	//fmt.Printf("\033[31m失败的数量：%v\033[0m \n", iLive-(oneLive+twoLive+threeLive)-okNumLive)
	//fmt.Printf("\033[31m总失败率：%.2f%v\033[0m \n", float64(iLive-(oneLive+twoLive+threeLive)-okNumLive)/float64(iLive-(oneLive+twoLive+threeLive))*100, "%")
	//fmt.Printf("最小响应时间：%.3f微秒 \n", float64(minRespTimeLive()))
	//fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTimeLive())/1e6)
	//fmt.Println("50%用户响应时间: " + fmt.Sprintf("%.3f秒", float64(fiftyRespTime())/1e6))
	//fmt.Println("90%用户响应时间: " + fmt.Sprintf("%.3f秒", float64(ninetyRespTime())/1e6))
	//fmt.Println(fmt.Sprintf("平均响应时间：%.3f秒", (float64(sumResTimeLive())/float64(len(resTimeLiveList)))/1e6))
	//fmt.Println(fmt.Sprintf("QPS: %.3f", float64(threeLive)/(float64(sumResTimeLive())/float64(len(resTimeLiveList))/1e6)))
	defer runtime.GC()
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
func liveSend(times int64, i int) {
	defer func() {
		err5 := recover()
		if err5 != nil {
			log.Println("捕获的异常：", err5)
		}
	}()
	//chMap := <-ch
	sTime = time.Now().UnixNano() / 1e6
	liveData := make(map[string]interface{})
	liveData["type"] = wsType
	liveData["client_name"] = userData[i]["nick_name"]
	liveData["room_id"] = roomId
	userId, _ := strconv.Atoi(userData[i]["id"])
	liveData["user_id"] = userId
	liveData["avatar"] = userData[i]["avatar"]
	liveData["token"] = userData[i]["token"]
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
	for {
		sendD := make(map[string]interface{})
		sendD["type"] = sendDataType
		sendD["content"] = contentlist[source.Intn(5)]
		sendBs, _ := json.Marshal(&sendD)
		s := time.Now().UnixNano() / 1e3
		err8 := ws.WriteMessage(websocket.BinaryMessage, sendBs)
		log.Println(`发送聊天信息：`, string(sendBs))
		if err8 != nil {
			log.Println("发送聊天信息异常：", err2)
		}
		atomic.AddInt64(&iLive, 1)
		_, recv, err3 := ws.ReadMessage()
		if err3 != nil {
			log.Println("接收数据异常", err3)
		}
		eTime = time.Now().UnixNano() / 1e6
		if eTime-sTime > times*1000 {
			break
		}
		var resMsg map[string]interface{}
		err4 := json.Unmarshal(recv, &resMsg)
		if err4 != nil {
			log.Println("解析返回数据异常：", err4)
		}
		log.Println(`服务端返回：`, string(recv))
		e := time.Now().UnixNano() / 1e3
		if string(recv) != "" {
			atomic.AddInt64(&okNumLive, 1)
		} else {
			fmt.Println("接收服务端返回数据异常：", string(recv))
		}
		lockLive.Lock()
		resTimeLiveList = append(resTimeLiveList, int(e-s))
		lockLive.Unlock()
	}
	//channelLive <- eTime - sTime
	defer ws.Close()
	defer wgLive.Done()
}
