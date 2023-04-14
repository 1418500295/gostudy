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
	"sort"
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
var reqUrl = "http"
var resTime int64
var chanResTime chan int64
var resTimeList []int
var rTimeChan = make(chan int, 100000) //创建响应数据收集缓冲区
var sucNum int64
var failNum int64
var lock sync.Mutex
var useTime int64
var avgSize int64

func avgResTime(timeList []int) float64 {
	sum := 0
	if len(timeList) == 0 {
		fmt.Println("列表数据为空")
	} else {
		for _, v := range timeList {
			sum += v
		}
	}
	return float64(sum) / float64(len(timeList))
}
func qps() float64 {
	return num / (avgResTime(resTimeList) / 1000)
}
func maxRespTime(timeList []int) float64 {
	max := timeList[0]
	for _, index := range timeList {
		if index > max {
			max = index
		}
	}
	return float64(max) / float64(1000)
}
func minRespTime(timeList []int) float64 {
	min := timeList[0]
	for _, index := range timeList {
		if index < min {
			min = index
		}
	}
	return float64(min) / float64(1000)
}
func fiftyRespTime(timeList []int) float64 {
	sort.Ints(timeList)
	resSize := 0.5
	return float64(timeList[int(float64(len(timeList))*resSize)-1]) / float64(1000)
}
func ninetyRespTime(timeList []int) float64 {
	sort.Ints(timeList)
	resSize := 0.9
	return float64(timeList[int(float64(len(timeList))*resSize)-1]) / float64(1000)
}

type Requests struct {
	url     string
	data    map[string]interface{}
	headers map[string]string
}

var reqData = map[string]interface{}{
	"cb":                  5,
	"
	//"start_work_start_time": time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).Unix(),
	//"start_work_end_time":   time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 0, time.Now().Location()).Unix(),
}
var headers = map[string]string{
	"t": "eyJhbGci",
}

//var transPort *http.Transport

//var req  *HttpRequest.Request
var cli = &http.Client{Transport: &http.Transport{
	MaxIdleConns:        10000, // Set your desired maximum number of idle connections
	MaxIdleConnsPerHost: 10000,
	IdleConnTimeout:     30 * time.Second, // Set your desired idle connection timeout
	DisableCompression:  true,             // Optional: Disable compression for testing purposes
	DisableKeepAlives:   false,            //复用连接
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
		avgSize = req.ContentLength
		rTimeChan <- int(e - s)
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
	table, err := gotable.Create("耗时", "并发数", "成功的请求数",
		"失败的请求数", "平均响应时间", "50%响应时间", "90%响应时间", "最大响应时间", "最小响应时间", "平均请求字节数", "qps")
	if err != nil {
		fmt.Println(err)
	}
	useTime = et - st
	table.AddRow([]string{fmt.Sprintf("%.1f", float64(useTime/1000)),
		strconv.Itoa(num), strconv.FormatInt(sucNum, 10),
		strconv.FormatInt(failNum, 10),
		fmt.Sprintf("%.3f", avgResTime(resTimeList)/1000),
		fmt.Sprintf("%.3f", fiftyRespTime(resTimeList)),
		fmt.Sprintf("%.3f", ninetyRespTime(resTimeList)),
		fmt.Sprintf("%.3f", maxRespTime(resTimeList)),
		fmt.Sprintf("%.3f", minRespTime(resTimeList)),
		strconv.FormatInt(avgSize, 10),
		fmt.Sprintf("%.3f", qps()),
	})
	fmt.Println(table)

}
func main() {
	//fmt.Println(utils.Md5ToStr("OK6802"))
	countDown := &sync.WaitGroup{}
	control := &sync.WaitGroup{}
	singleChan := make(chan struct{})
	//account1.StartWork()
	this := Requests{headers: headers, url: reqUrl, data: reqData}
	this.run(5, 10, countDown, control, singleChan, 15)
	printTable()
}
