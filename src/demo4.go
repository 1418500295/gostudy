package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/liushuochen/gotable"
	"github.com/marusama/cyclicbarrier"
	"io/ioutil"
	"net/http"
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

// var token3 = "eyJhbGc"
var reqUrl = "https://test-www.vipsroom.net/api/scene/convert/chips/start/work/list"
var resTime int64
var chanResTime chan int64
var resTimeList []int64
var rTimeChan = make(chan int64, 100000) //创建响应数据收集缓冲区
var sucNum int64
var failNum int64
var lock sync.Mutex
var useTime int64

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
	"1":                  5,
	
	
}
var headers = map[string]string{
	"t": "",
}

//var transPort *http.Transport

//var req  *HttpRequest.Request
var cli = &http.Client{Transport: &http.Transport{
	MaxIdleConns:        10000, // Set your desired maximum number of idle connections
	MaxIdleConnsPerHost: 1000,
	IdleConnTimeout:     30 * time.Second, // Set your desired idle connection timeout
	DisableCompression:  true,             // Optional: Disable compression for testing purposes
}}

func getData() *bytes.Reader {
	dataBy, _ := json.Marshal(reqData)
	reader := bytes.NewReader(dataBy)
	return reader
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
	st = time.Now().UnixMilli()
	for true {
		//req := HttpRequest.NewRequest()
		//headers := map[string]string{"t": token3}
		var err error
		req, err := http.NewRequest("POST", requests.url, getData())
		req.Header.Add("content-type", "application/json")
		req.Header.Add("t", headers["t"])
		s := time.Now().UnixNano() / 1e6
		res, _ := cli.Do(req)
		e := time.Now().UnixNano() / 1e6
		body, _ := ioutil.ReadAll(res.Body)
		rTimeChan <- e - s
		//lock.Lock()
		//resTimeList = append(resTimeList, resTime)
		//lock.Unlock()
		//fmt.Printf("第%d协程请求返回:%s\n", i, string(body))
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
			//defer res.Close()
			et = time.Now().UnixMilli()
			if et-st > times*1000 {
				break
			}
			continue
		}

	}
	//chanResTime <- resTime
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
			err := barrier.Await(context.Background()) //同步集合点
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
		countDown.Wait() //等待所有子协程执行完毕
		singleChan <- struct{}{}
		close(rTimeChan) //关闭数据收集通道
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
	//fmt.Println("------执行结束------")
	//fmt.Println("成功的请求数：", sucNum)
	//fmt.Println("失败的请求数：", failNum)
	//fmt.Printf("平均响应时间：%.3f毫秒", avgResTime(resTimeList))

}

func printTable() {
	table, err := gotable.Create("耗时", "总并发数", "成功的请求数", "失败的请求数", "平均响应时间")
	if err != nil {
		fmt.Println(err)
	}
	useTime = et - st
	table.AddRow([]string{fmt.Sprintf("%.1f", float64(useTime/1000)),
		strconv.Itoa(num), strconv.FormatInt(sucNum, 10),
		strconv.FormatInt(failNum, 10),
		fmt.Sprintf("%.3f", avgResTime(resTimeList)/1000),
	})
	fmt.Println(table)

}
func main() {
	//fmt.Println(utils.Md5ToStr(""))
	countDown := &sync.WaitGroup{}
	control := &sync.WaitGroup{}
	singleChan := make(chan struct{})
	//account1.StartWork()
	this := Requests{headers: headers, url: reqUrl, data: reqData}
	this.run(5, 10, countDown, control, singleChan, 15)
	printTable()
}
