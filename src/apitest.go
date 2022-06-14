package main

import (
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	//创建计数器
	wg             = sync.WaitGroup{}
	num      int64 = 500 //设置并发数量
	okNum    int64 = 0   //初始化请求成功的数量
	timeList []int       //响应时间
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

func main() {
	//格式化输出时间
	//start_time := time.Now().Format("2006-01-02 15:04:05")
	startTime := time.Now().UnixNano() / 1e6
	fmt.Printf("开始时间：%v \n", startTime)
	do(num)
	endTime := time.Now().UnixNano() / 1e6
	fmt.Printf("结束时间：%v \n", endTime)
	fmt.Println("成功的数量：", okNum)
	fmt.Printf("失败的数量：%v \n", num-okNum)
	fmt.Printf("总耗时：%.3f 秒 \n", float64(endTime-startTime)/1000)
	fmt.Printf("最大响应时间：%.3f 秒 \n", float64(maxRespTime())/1000)
	fmt.Printf("最小响应时间：%.3f 秒 \n", float64(minRespTime())/1000)
	fmt.Printf("平均响应时间是:%.3f 秒 \n", float64(sumRespTime())/float64(num)/1000)
	fmt.Printf("QPS：%.3f", float64(num)/(float64(sumRespTime())/float64(num)/1000))
	//确保打包成exe运行后窗口不关闭
	_, _ = fmt.Scanf("h")
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
	request.url = ""
	data := make(map[string]interface{})
	data["userName"] = "admin"
	data["password"] = "111"
	data["safeCode"] = "121"
	sTime := time.Now().UnixNano() / 1e6
	resp, _ := HttpRequest.Post(request.url, data)
	defer resp.Close()
	if resp.StatusCode() == 200 {
		body, _ := resp.Body()
		fmt.Println(string(body))
		eTime := time.Now().UnixNano() / 1e6
		use_time := (eTime - sTime)
		channel <- use_time
		//time_list = append(time_list, int(use_time))
		jsonData := gjson.Parse(string(body)).Map()
		if jsonData["code"].Int() == -1 {
			// 多个goroutine并发读写sum，有并发冲突，最终计算得到的sum值是不准确的
			// 使用原子操作计算sum，没有并发冲突，最终计算得到sum的值是准确的
			atomic.AddInt64(&okNum, 1)
		}
	} else {
		panic("请求异常")
	}
	//计数器减一
	defer wg.Done()

}
