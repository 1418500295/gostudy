package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"gostudy/src/protouse"
	"log"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	time "time"
	"unsafe"
)

const wsUri = "ws://sports.aisport.live/ws?token=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjE4NTY0NzgsImlhdCI6MTY2MTQ5NjQ3OCwiaXNzIjoic3BvcnRib29rX2FwaSIsInN1YiI6InNwb3J0IiwiT3BlcmF0b3JJZCI6MjA5LCJTaXRlSWQiOjY0LCJVc2VySWQiOjEyODk5NDYzOTYsIkN1cnJlbmN5SWQiOjQsIkFjY291bnQiOiJlZGVuMTExMSJ9.8QMT22SKyBhvEY8sD7hVgeiK_UGurw70HCvrhGuPA7lcjUVQfHPFoP0PNxOSE-DQd-vo4sJCftbKUYoAqAXRpA"

//var origin = "http://127.0.0.1:8080/"
//var url = "ws://127.0.0.1:8080/echo"
var (
	wg = sync.WaitGroup{}
	//num         int64 = 2
	okNum       int64 = 0
	channel1          = make(chan int64)
	resTimeList []int
	lock1             = sync.Mutex{}
	i1          int64 = 0
	done1             = make(chan struct{})
	//resTime     int64
	one            int64
	two            int64
	three          int64
	first          int64
	second         int64
	third          int64
	connectNum     int64 = 0
	connectFailNum int64 = 0

	iFirst       int64
	iSecond      int64
	iThird       int64
	okFirst      int64
	okSecond     int64
	okThird      int64
	resTimelist1 []int
	resTimelist2 []int
	resTimelist3 []int
	starts       int64
	ends         int64
	useT1        int64
	useT2        int64
	useT3        int64
)

func maxRespTime(resTimeList []int) int {
	max := resTimeList[0]
	for _, i := range resTimeList {
		if i > max {
			max = i
		}
	}
	return max
}
func minRespTime(resTimeList []int) int {
	min := resTimeList[0]
	for _, i := range resTimeList {
		if i < min {
			min = i
		}
	}
	return min
}
func sumResTime(resTimeList []int) int {
	sum := 0
	for _, i := range resTimeList {
		sum += i
	}
	return sum
}

func fiftyTime(resTimeList []int) int {
	sort.Ints(resTimeList)
	resSize := 0.5
	return resTimeList[int(float64(len(resTimeList))*resSize)-1]
}
func ninetyTime(resTimeList []int) int {
	sort.Ints(resTimeList)
	resSize := 0.9
	return resTimeList[int(float64(len(resTimeList))*resSize)-1]
}

func printTime(usetimes int64, count int64, i1 int64, okNumLive int64, resTimeList []int) {
	fmt.Println("并发数：", count)
	fmt.Println("请求数: ", i1)
	fmt.Println("成功的数量：", len(resTimeList))
	fmt.Printf("\033[31m失败的数量：%v\033[0m \n", int(i1)-len(resTimeList))
	fmt.Printf("\033[31m失败率：%.2f%v\033[0m \n", float64(int(i1)-len(resTimeList))/float64(i1)*100, "%")
	fmt.Println(fmt.Sprintf("耗时: %v秒", float64(usetimes)/1000))
	fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTime(resTimeList))/1e6)
	fmt.Printf("最小响应时间：%.3f微秒 ≈ %v秒 \n", float64(minRespTime(resTimeList)), float64(minRespTime(resTimeList))/1e6)
	fmt.Println("50%用户响应时间: " + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(fiftyTime(resTimeList)), float64(fiftyTime(resTimeList))/1e6))
	fmt.Println("90%用户响应时间: " + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(ninetyTime(resTimeList)), float64(ninetyTime(resTimeList))/1e6))
	fmt.Printf("平均响应时间是:%.3f秒 \n", float64(sumResTime(resTimeList))/float64(i1)/1e6)
	fmt.Printf("QPS：%.3f \n", float64(count)/(float64(sumResTime(resTimeList))/float64(i1)/1e6))
}
func run() {
	fmt.Println(fmt.Sprintf("***首轮并发用户为%v***", one))
	wg.Add(int(one))
	for i := 0; i < int(one); i++ {
		go send(first)
	}
	wg.Wait()
	iFirst = i1
	okFirst = okNum
	resTimelist1 = resTimeList
	useT1 = ends - starts
	fmt.Println("第二轮压测开始...")
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", two))
	<-time.After(10 * time.Second)
	wg.Add(int(two))
	for i := 0; i < int(two); i++ {
		go send(second)
	}
	wg.Wait()
	iSecond = i1 - iFirst
	okSecond = okNum - okFirst
	resTimelist2 = resTimeList[len(resTimelist1):]
	useT2 = ends - starts
	fmt.Println("第三轮压测开始...")
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", three))
	<-time.After(10 * time.Second)
	wg.Add(int(three))
	for i := 0; i < int(three); i++ {
		go send(third)
	}
	wg.Wait()
	iThird = i1 - iFirst - iSecond
	okThird = okNum - okFirst - okSecond
	resTimelist3 = resTimeList[len(resTimelist1)+len(resTimelist2):]
	useT3 = ends - starts
	fmt.Println("\033[32m建立连接成功的数量：\033[0m", connectNum)
	fmt.Println("\033[31m建立连接失败的数量：\033[0m", connectFailNum)
	fmt.Println("\033[33m***第一轮压测结果***\033[0m")
	printTime(useT1, one, iFirst, okFirst, resTimelist1)
	fmt.Println("\033[35m----------------------- \033[0m")
	fmt.Println("\033[33m***第二轮压测结果***\033[0m")
	printTime(useT2, two, iSecond, okSecond, resTimelist2)
	fmt.Println("\033[35m----------------------- \033[0m")
	fmt.Println("\033[33m***第三轮压测结果***\033[0m")
	printTime(useT3, three, iThird, okThird, resTimelist3)
	fmt.Println("\033[35m----------------------- \033[0m")
}
func main() {
	//flag.Int64Var(&one, "oneCount", 10000, "第一轮并发数")
	//flag.Int64Var(&first, "oneTime", 100, "第一轮压测运行时长(秒)")
	//flag.Int64Var(&two, "twoCount", 15000, "第二轮并发数")
	//flag.Int64Var(&second, "twoTime", 100, "第二轮压测运行时长(秒)")
	//flag.Int64Var(&three, "threeCount", 20000, "第三轮并发数")
	//flag.Int64Var(&third, "threeTime", 100, "第三轮压测运行时长(秒)")
	//flag.Parse()
	fmt.Println("第一轮压测用户数：")
	_, _ = fmt.Scan(&one)
	fmt.Println("第一轮运行时长(秒)：")
	_, _ = fmt.Scan(&first)
	fmt.Println("第二轮压测用户数：")
	_, _ = fmt.Scan(&two)
	fmt.Println("第二轮运行时长(秒): ")
	_, _ = fmt.Scan(&second)
	fmt.Println("第三轮压测用户数：")
	_, _ = fmt.Scan(&three)
	fmt.Println("第三轮运行时长(秒): ")
	_, _ = fmt.Scan(&third)
	sTimes := time.Now().UnixNano() / 1e6
	run()
	eTimes := time.Now().UnixNano() / 1e6
	fmt.Println(fmt.Sprintf("总耗时：%v秒", float64(eTimes-sTimes)/1000-(5+5)))
	defer runtime.GC()
}

var payLoad protouse.PayloadAction

func init() {
	payLoad = protouse.PayloadAction{
		Action: 0, //subscribing
		MatchIds: []string{
			"3552578700",
			"3557513700",
		},
		MarketTypeIds: []uint64{
			1,
			16,
			18,
			60,
			66,
			68,
			16,
			18,
			66,
			68,
		},
		Scope: 1,
	}
}

func send(times int64) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	byteParam, err := proto.Marshal(&payLoad)
	if err != nil {
		log.Println("\003[31m解析请求参数异常：\033[0m", err)
	}
	ws, con, err1 := websocket.DefaultDialer.Dial(wsUri, nil)
	//fmt.Println(con.Header)
	if err1 != nil {
		log.Println("\033[31m建立连接异常: \033[0m", err1)
		atomic.AddInt64(&connectFailNum, 1)
	}
	if con.StatusCode == 101 {
		atomic.AddInt64(&connectNum, 1)
		//log.Println("\033[32m已成功建立的连接数：\033[0m", connectNum)
	}
	starts = time.Now().UnixNano() / 1e6
	for {
		s := time.Now().UnixNano() / 1e6
		err2 := ws.WriteMessage(websocket.BinaryMessage, byteParam)
		//log.Printf("客户端发送：%v\n", payLoad)
		if err2 != nil {
			log.Println("\033[31m客户端发送数据异常：\033[0m", err2)
		}
		atomic.AddInt64(&i1, 1)
		_, recv, err3 := ws.ReadMessage()
		if err3 != nil {
			log.Println("\033[31m接收服务端返回数据异常\033[0m", err3)
		} else {
			atomic.AddInt64(&okNum, 1)
		}
		e := time.Now().UnixNano() / 1e6
		ends = time.Now().UnixNano() / 1e6
		log.Printf("响应时间：%v毫秒", e-s)
		var resMsg protouse.Message
		err4 := proto.Unmarshal(recv, &resMsg)
		if err4 != nil {
			log.Println("\033[31m解析返回数据异常：\033[0m", err4)
		}
		//log.Printf("服务端返回：%v", resMsg.Data.Markets[0])
		lock1.Lock()
		resTimeList = append(resTimeList, int((e-s)*1000))
		lock1.Unlock()
		if ends-starts > times*1000 {
			break
		}
		//ws.Close()
	}

}
