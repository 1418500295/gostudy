package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"gostudy/src/protouse"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	time "time"
)

const wsUri = "ws://sports.aisport.live/ws?token=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTc5OTI5NDksImlhdCI6MTY1NzYzMjk0OSwiaXNzIjoic3BvcnRib29rX2FwaSIsInN1YiI6InNwb3J0IiwiT3BlcmF0b3JJZCI6MjA5LCJTaXRlSWQiOjY0LCJVc2VySWQiOjEyODk5NDYzOTYsIkN1cnJlbmN5SWQiOjQsIkFjY291bnQiOiJlZGVuMTExMSJ9.LtfdtjD1bOGxnW6CW3A8XtORhy7sB5WtWic1ByB6cwy13COAF5323j_gUAh3H4LxXLNYLoCujZW7rKELRYpv5g"

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
	one        int64
	two        int64
	three      int64
	first      int64
	second     int64
	third      int64
	connectNum int64 = 0
)

func maxRespTime() int {
	max := resTimeList[0]
	for _, i := range resTimeList {
		if i > max {
			max = i
		}
	}
	return max
}
func minRespTime() int {
	min := resTimeList[0]
	for _, i := range resTimeList {
		if i < min {
			min = i
		}
	}
	return min
}
func sumResTime() int {
	sum := 0
	for _, i := range resTimeList {
		sum += i
	}
	return sum
}

func printTime(usetime int64, count int64) {
	fmt.Println("并发数：", count)
	fmt.Println("请求数: ", i1)
	fmt.Println("成功的数量：", okNum)
	fmt.Printf("失败的数量：%v \n", i1-okNum)
	fmt.Println(fmt.Sprintf("耗时: %v秒", float64(usetime)/1000))
	fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTime())/1000)
	fmt.Printf("最小响应时间：%.3f毫秒 \n", float64(minRespTime()))
	fmt.Printf("平均响应时间是:%.3f秒 \n", float64(sumResTime())/float64(i1)/1000)
	fmt.Printf("QPS：%.3f \n", float64(count)/(float64(sumResTime())/float64(i1)/1000))
}
func run() {
	//wg.Add(int(count))
	//for i := 0; i < int(count); i++ {
	//	go send(times)
	//	//go func() {
	//	//	resTime := <-channel1
	//	//	resTimeList = append(resTimeList, int(resTime))
	//	//}()
	//}
	fmt.Println(fmt.Sprintf("***首轮并发用户为%v***", one))
	fsT := time.Now().UnixNano() / 1e6
	wg.Add(int(one))
	for i := 0; i < int(one); i++ {
		go send(first)
	}
	wg.Wait()
	feT := time.Now().UnixNano() / 1e6
	fmt.Println("***第一轮压测结果***")
	printTime(feT-fsT, one)
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", two))
	<-time.After(10 * time.Second)
	seSt := time.Now().UnixNano() / 1e6
	wg.Add(int(two))
	for i := 0; i < int(two); i++ {
		go send(second)
	}
	wg.Wait()
	seEd := time.Now().UnixNano() / 1e6
	fmt.Println("***第二轮压测结果***")
	printTime(seEd-seSt, two)
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", three))
	<-time.After(10 * time.Second)
	wg.Add(int(three))
	for i := 0; i < int(three); i++ {
		go send(third)
	}
	wg.Wait()
	//go func() {
	//	wg.Wait()
	//	close(done1)
	//}()
	//select {
	//case <-done1:
	//	fmt.Println("任务完成")
	//case <-time.After(50 * time.Minute):
	//	panic("任务处理超时")
	//}
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
	sTime := time.Now().UnixNano() / 1e6
	run()
	eTime := time.Now().UnixNano() / 1e6
	fmt.Println(fmt.Sprintf("总耗时：%v秒", float64(eTime-sTime)/1000-(5+5)))
	fmt.Println("总并发数：", one+two+three)
	fmt.Println("总请求数: ", i1)
	fmt.Printf("成功的数量：%v \n", okNum)
	fmt.Printf("失败的数量：%v \n", i1-okNum)
	fmt.Printf("最小响应时间：%.3f毫秒 \n", float64(minRespTime()))
	fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTime())/1000)
	fmt.Println(fmt.Sprintf("平均响应时间：%.3f秒", (float64(sumResTime())/float64(len(resTimeList)))/1000))
	fmt.Println(fmt.Sprintf("QPS: %.3f", float64(one+two+three)/(float64(sumResTime())/float64(len(resTimeList))/1000)))
	defer runtime.GC()
}

func send(times int64) {
	defer func() {
		err5 := recover()
		if err5 != nil {
			log.Println("捕获的异常：", err5)
		}
	}()
	sT := time.Now().UnixNano() / 1e6
	payLoad := protouse.PayloadAction{
		Action: 0, //subscribing
		MatchIds: []string{
			"3189041300",
			"3297110700",
			"3374500500",
			"3422906500",
			"3439196700",
			"3441710900",
			"3441711300",
			"3441711900",
			"3441712100",
			"3453772700",
		},
		MarketTypeIds: []uint64{1, 16, 18, 60, 66, 68, 16, 1182, 66, 68},
		Scope:         1,
	}
	byteParam, err := proto.Marshal(&payLoad)
	if err != nil {
		log.Println("解析请求参数异常：", err)
	}
	bs, _ := json.Marshal(&payLoad)
	for {
		ws, res, err1 := websocket.DefaultDialer.Dial(wsUri, nil)
		if err1 != nil {
			log.Println(fmt.Sprintf("建立连接异常：%v", err1))
		}
		if res != nil {
			log.Println(`连接成功: `, res)
			atomic.AddInt64(&connectNum, 1)
		}
		log.Println("已成功建立的连接数：", connectNum)
		sTime := time.Now().UnixNano() / 1e6
		err2 := ws.WriteMessage(websocket.BinaryMessage, byteParam)
		log.Println(`客户端发送：`, string(bs))
		atomic.AddInt64(&i1, 1)
		if err2 != nil {
			log.Println("发送数据异常：", err2)
		}
		_, recv, err3 := ws.ReadMessage()
		if err3 != nil {
			log.Println("接收数据异常", err3)
		}
		eT := time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		var resMsg protouse.Message
		err4 := proto.Unmarshal(recv, &resMsg)
		if err4 != nil {
			log.Println("解析返回数据异常：", err4)
		}
		eTime := time.Now().UnixNano() / 1e6
		lock1.Lock()
		resTimeList = append(resTimeList, int(eTime-sTime))
		lock1.Unlock()
		byS, _ := json.Marshal(resMsg)
		log.Println(`服务端返回：`, string(byS))
		if string(byS) != "" {
			atomic.AddInt64(&okNum, 1)
		} else {
			fmt.Println("响应异常：", string(byS))
		}
	}
	//channel1 <- eTime - sTime
	defer wg.Done()
}
