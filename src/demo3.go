package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"github.com/marusama/cyclicbarrier"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const num = 5

var et int64
var st int64
var token3 = "..5yx8ENOZ02yKJCgq2Ltrp0zpknQpRzI6yS2wAAzcDgE"
var reqUrl = ""

var resTime int64
var chanResTime chan int64
var resTimeList []int64
var rTimeChan = make(chan int64, 100000)
var sucNum int64
var failNum int64
var lock sync.Mutex

func avgResTime(timeList []int64) float64 {
	sum := 0
	if len(timeList) == 0 {
		fmt.Println("列表数据为空")
	} else {
		for _, v := range timeList {
			sum += int(v)
		}
	}
	return float64(sum) / float64(len(timeList))
}

type Requests struct {
	url     string
	data    map[string]interface{}
	headers map[string]string
}

var reqData = map[string]interface{}{
	"member_code":           "",
	"page":                  1,
	"size":                  10,
	"start_work_start_time": time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).Unix(),
	"start_work_end_time":   time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 0, time.Now().Location()).Unix(),
}
var headers = map[string]string{
	"t": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVC.",
}

func (requests *Requests) execute(i int, times int64) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	st = time.Now().Unix()
	for true {
		req := HttpRequest.NewRequest()
		//headers := map[string]string{"t": token3}
		var err error
		res, err := req.SetHeaders(requests.headers).JSON().Post(requests.url, requests.data)
		if err != nil {
			fmt.Printf("请求异常：%v\n", err)
		} else {
			var rTime string
			resT := func() string {
				if len(res.Time()) != 0 {
					rTime = res.Time()[:len(res.Time())-2]
				}
				return rTime
			}
			resTime, err = strconv.ParseInt(resT(), 10, 64)
			if err != nil {
				fmt.Printf("响应时间转换异常：%v\n", err)
			} else {
				rTimeChan <- resTime
				//lock.Lock()
				//resTimeList = append(resTimeList, resTime)
				//lock.Unlock()
				body, _ := res.Body()
				fmt.Printf("第%d协程请求返回:%s\n", i, string(body))
				var resMap map[string]interface{}
				err = json.Unmarshal(body, &resMap)
				if err != nil {
					fmt.Printf("解析返回数据异常：%v\n", err)
				} else {
					if resMap["msg"].(string) == `请求成功` && resMap["code"].(float64) == 10000 {
						atomic.AddInt64(&sucNum, 1)
					} else {
						atomic.AddInt64(&failNum, 1)
					}
					defer res.Close()
					et = time.Now().Unix()
					if et-st > times {
						break
					}
					continue
				}
			}
		}
		//chanResTime <- resTime
	}
}

func (requests *Requests) run(num int, executeTimes int64, countDown *sync.WaitGroup, control *sync.WaitGroup, singleChan chan struct{}, timeOut int) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("safe go: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	barrier := cyclicbarrier.New(num)
	countDown.Add(num)
	control.Add(1)
	for i := 0; i < num; i++ {
		go func(i int) {
			control.Wait() //开启阀门，阻塞子协程执行
			requests.execute(i, executeTimes)
			err := barrier.Await(context.Background()) //同步屏障触发点
			if err != nil {
				return
			}
			defer countDown.Done() //计数器-1
		}(i)
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancle()
	fmt.Println("------执行开始------")
	control.Done() //关闭阀门
	go func() {
		countDown.Wait() //所有子协程执行完毕
		singleChan <- struct{}{}
		close(rTimeChan)
	}()
	select {
	case <-singleChan:
		fmt.Println("******协程处理完成******")
	case <-ctx.Done():
		fmt.Println("协程处理超时!!!")
	}
	for ch := range rTimeChan {
		resTimeList = append(resTimeList, ch)
	}
	fmt.Println("------执行结束------")
	fmt.Println("成功的请求数：", sucNum)
	fmt.Println("失败的请求数：", failNum)
	fmt.Printf("平均响应时间：%.3f毫秒", avgResTime(resTimeList))

}
func main() {
	countDown := &sync.WaitGroup{}
	control := &sync.WaitGroup{}
	singleChan := make(chan struct{})
	//account1.StartWork()
	this := Requests{headers: headers, url: reqUrl, data: reqData}
	this.run(3, 5, countDown, control, singleChan, 10)
}
