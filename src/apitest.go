package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	//创建计数器
	wg             = sync.WaitGroup{}
	num      int64 = 10 //设置并发数量
	okNum    int64 = 0  //初始化请求成功的数量
	timeList []int      //响应时间
	channel  = make(chan int64)
)

type Request struct {
	url  string
	data map[string]interface{}
}

//获取时间戳
func SetTs() string {
	return strconv.FormatInt(time.Now().Unix()*1000, 10)
}

func sumRespTime() int {
	sum := 0
	for _, i := range timeList {
		sum = sum + i
	}
	return sum
}

func maxRespTime() int {
	max := timeList[0]
	for _, i := range timeList {
		if i > max {
			max = i
		}
	}
	return max
}
func minRespTime() int {
	min := timeList[0]
	for _, i := range timeList {
		if i < min {
			min = i
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

func main() {
	//格式化输出时间
	//start_time := time.Now().Format("2006-01-02 15:04:05")
	startTime := time.Now().UnixNano() / 1e6
	fmt.Printf("开始时间：%v \n", startTime)
	do(num)
	endTime := time.Now().UnixNano() / 1e6
	fmt.Printf("结束时间：%v \n", endTime)
	fmt.Println("总请求数: ", num)
	fmt.Println("成功的数量：", okNum)
	fmt.Printf("失败的数量：%v \n", num-okNum)
	fmt.Printf("总耗时：%.3f 秒 \n", float64(endTime-startTime)/1000)
	fmt.Println("50%用户响应时间：" + fmt.Sprintf("%.3f秒", float64(fiftyRespTime())/1000))
	fmt.Println("90%用户响应时间：" + fmt.Sprintf("%.3f秒", float64(ninetyRespTime())/1000))
	fmt.Printf("最大响应时间：%.3f 秒 \n", float64(maxRespTime())/1000)
	fmt.Printf("最小响应时间：%.3f 秒 \n", float64(minRespTime())/1000)
	fmt.Printf("平均响应时间是:%.3f 秒 \n", float64(sumRespTime())/float64(num)/1000)
	fmt.Printf("QPS：%.3f", float64(num)/(float64(sumRespTime())/float64(num)/1000))
	runtime.GC()
	//确保打包成exe运行后窗口不关闭
	//_, _ = fmt.Scanf("h")
}

func do(num int64) {
	//设置多核CPU运行
	//runtime.GOMAXPROCS(runtime.NumCPU())
	wg.Add(int(num)) //初始化计数器
	for i := 0; i < int(num); i++ {
		go httpSend()
		//解决高并发下Slice写入的线程安全问题
		go func() {
			data := <-channel
			timeList = append(timeList, int(data))
		}()
	}
	//主协程阻塞子协程执行完毕
	wg.Wait()
	defer close(channel)
}

//发送请求
func httpSend() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("错误信息: ", err)
		}
	}()
	var request Request
	request.url = "http://192.168.128.156:3333/boss/login"
	//data := make(map[string]interface{})
	//data["userName"] = ""
	//data["password"] = ""
	//data["safeCode"] = ""
	//data := sync.Map{}
	//data.Store("userName", "admin")
	//data.Store("password", "111")
	//data.Store("safeCode", "121")
	//jsByte, _ := json.Marshal(&data)
	//var jsonData map[string]interface{}
	//err := json.Unmarshal(jsByte, &jsonData)
	//if err != nil {
	//	fmt.Println(err)
	//}
	args := &fasthttp.Args{}
	args.Add("userName", "admin")
	args.Add("password", "111")
	args.Add("safeCode", "121")
	sTime := time.Now().UnixNano() / 1e6
	c := &fasthttp.Client{
		MaxConnsPerHost: 10000,
		ReadTimeout:     4000 * time.Millisecond,
		WriteTimeout:    4000 * time.Millisecond,
	}
	_, res, err := c.Post(nil, request.url, args)
	eTime := time.Now().UnixNano() / 1e6
	channel <- eTime - sTime
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
	if gjson.ParseBytes(res).Map()["code"].Int() == -1 {
		atomic.AddInt64(&okNum, 1)
	}

	//resp, _ := HttpRequest.Post(request.url, data)
	//defer resp.Close()
	//if resp.StatusCode() == 200 {
	//	body, _ := resp.Body()
	//	fmt.Println(string(body))
	//	respTime, _ := strconv.ParseInt(strings.Split(resp.Time(), "m")[0], 10, 64)
	//	channel <- respTime
	//	jsonData := gjson.Parse(string(body)).Map()
	//	if jsonData["code"].Int() == -1 {
	//		// 多个goroutine并发读写sum，有并发冲突，最终计算得到的sum值是不准确的
	//		// 使用原子操作计算sum，没有并发冲突，最终计算得到sum的值是准确的
	//		atomic.AddInt64(&okNum, 1)
	//	}
	//} else {
	//	panic("请求异常")
	//}
	//计数器减一
	defer wg.Done()

}
