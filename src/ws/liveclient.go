package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
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
	wgLive = sync.WaitGroup{}
	//num         int64 = 2
	okNumLive        int64 = 0
	firstOkNumLive   int64
	secondOkNumLive  int64
	thirdOkNumLive   int64
	channelLive      = make(chan int64)
	resTimeLiveList  []int
	resTimeLiveList1 []int
	resTimeLiveList2 []int
	resTimeLiveList3 []int
	lockLive               = sync.Mutex{}
	iLive            int64 = 0
	iLiveFirst       int64
	iLiveSecond      int64
	iLiveThird       int64
	doneLive1        = make(chan struct{})
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

	hashChan = make(chan string)
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
	fmt.Println("并发数：", count)
	fmt.Println("请求数: ", iLive-count)
	fmt.Println("成功的数量：", okNumLive)
	fmt.Printf("失败的数量：%v \n", iLive-count-okNumLive)
	fmt.Println(fmt.Sprintf("耗时: %v秒", float64(usetime)/1e6))
	fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTimeLive(resTimeLiveList))/1e6)
	fmt.Printf("最小响应时间：%.3f微秒 ≈ %v秒 \n", float64(minRespTimeLive(resTimeLiveList)), float64(minRespTimeLive(resTimeLiveList))/1e6)
	fmt.Println("50%用户响应时间: " + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(fiftyRespTime(resTimeLiveList)), float64(fiftyRespTime(resTimeLiveList))/1e6))
	fmt.Println("90%用户响应时间: " + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(ninetyRespTime(resTimeLiveList)), float64(ninetyRespTime(resTimeLiveList))/1e6))
	fmt.Printf("平均响应时间是:%.3f秒 \n", float64(sumResTimeLive(resTimeLiveList))/float64(iLive)/1e6)
	fmt.Printf("QPS：%.3f \n", float64(count)/(float64(sumResTimeLive(resTimeLiveList))/float64(iLive)/1e6))
}
func runLive() {
	fmt.Println(fmt.Sprintf("***首轮并发用户为%v***", oneLive))
	fsT := time.Now().UnixNano() / 1e6
	wgLive.Add(int(oneLive))
	for i := 0; i < int(oneLive); i++ {
		go liveSend(firstLive, i)
	}
	wgLive.Wait()
	feT := time.Now().UnixNano() / 1e6
	iLiveFirst = iLive
	firstOkNumLive = okNumLive
	resTimeLiveList1 = resTimeLiveList
	fmt.Println("***第一轮压测结果***")
	printTimeLive(feT-fsT, oneLive, iLiveFirst, firstOkNumLive, resTimeLiveList1)
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", twoLive))
	<-time.After(10 * time.Second)
	seSt := time.Now().UnixNano() / 1e6
	wgLive.Add(int(twoLive))
	for i := 0; i < int(twoLive); i++ {
		go liveSend(secondLive, i)
	}
	wgLive.Wait()
	seEd := time.Now().UnixNano() / 1e6
	iLiveSecond = iLive - iLiveFirst
	secondOkNumLive = okNumLive - firstOkNumLive
	resTimeLiveList2 = resTimeLiveList[len(resTimeLiveList1):]
	fmt.Println("***第二轮压测结果***")
	printTimeLive(seEd-seSt, twoLive, iLiveSecond, secondOkNumLive, resTimeLiveList2)
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", threeLive))
	<-time.After(10 * time.Second)
	thSt := time.Now().UnixNano() / 1e6
	wgLive.Add(int(threeLive))
	for i := 0; i < int(threeLive); i++ {
		go liveSend(thirdLive, i)
	}
	wgLive.Wait()
	thEt := time.Now().UnixNano() / 1e6
	iLiveThird = iLive - iLiveFirst - iLiveSecond
	thirdOkNumLive = okNumLive - firstOkNumLive - secondOkNumLive
	resTimeLiveList3 = resTimeLiveList[len(resTimeLiveList1)+len(resTimeLiveList2):]
	fmt.Println("***第三轮压测结果***")
	printTimeLive(thEt-thSt, threeLive, iLiveThird, thirdOkNumLive, resTimeLiveList3)
}

func main() {

	//byteData, err := ioutil.ReadFile(path)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err1 := json.Unmarshal(byteData, &userData)
	//if err1 != nil {
	//	fmt.Println(err1)
	//}
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

	byteData1, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	err1 := json.Unmarshal(byteData1, &userData)
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
	sTime := time.Now().UnixNano() / 1e6
	runLive()
	eTime := time.Now().UnixNano() / 1e6
	fmt.Println(fmt.Sprintf("总耗时：%v秒", float64(eTime-sTime)/1000-(10+10)))
	fmt.Println("总并发数：", oneLive+twoLive+threeLive)
	fmt.Println("总请求数: ", iLive-(oneLive+twoLive+threeLive))
	//fmt.Printf("成功的数量：%v \n", okNumLive)
	//fmt.Printf("失败的数量：%v \n", iLive-(oneLive+twoLive+threeLive)-okNumLive)
	//fmt.Printf("最小响应时间：%.3f微秒 \n", float64(minRespTimeLive()))
	//fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTimeLive())/1e6)
	//fmt.Println("50%用户响应时间: " + fmt.Sprintf("%.3f秒", float64(fiftyRespTime())/1e6))
	//fmt.Println("90%用户响应时间: " + fmt.Sprintf("%.3f秒", float64(ninetyRespTime())/1e6))
	//fmt.Println(fmt.Sprintf("平均响应时间：%.3f秒", (float64(sumResTimeLive())/float64(len(resTimeLiveList)))/1e6))
	//fmt.Println(fmt.Sprintf("QPS: %.3f", float64(threeLive)/(float64(sumResTimeLive())/float64(len(resTimeLiveList))/1e6)))
	defer runtime.GC()
}

func liveSend(times int64, i int) {
	defer func() {
		err5 := recover()
		if err5 != nil {
			log.Println("捕获的异常：", err5)
		}
	}()
	sT := time.Now().UnixNano() / 1e6
	liveData := make(map[string]interface{})
	liveData["type"] = wsType
	liveData["client_name"] = userData[i]["nick_name"]
	liveData["room_id"] = roomId

	userId, _ := strconv.Atoi(userData[i]["id"])
	liveData["user_id"] = userId
	fmt.Println(userId)
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
		sTime := time.Now().UnixNano() / 1e3
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
		eT := time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
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
		lockLive.Lock()
		resTimeLiveList = append(resTimeLiveList, int(eTime-sTime))
		lockLive.Unlock()
	}
	//channelLive <- eTime - sTime
	defer ws.Close()
	defer wgLive.Done()
}
