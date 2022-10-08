package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gostudy/src/protouse"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	token = "1"
)

var (
	//创建计数器
	wg = sync.WaitGroup{}
	//num      int64 = 5 //设置并发数量
	sT          int64
	eT          int64
	useTime1    int64
	useTime2    int64
	useTime3    int64
	okNum       int64 = 0 //初始化请求成功的数量
	firstOkNum  int64     //第一轮成功请求数
	secondOkNum int64     //第二轮成功请求数
	thirdOkNum  int64     //第三轮成功请求数
	timeList1   []int     //第一轮响应时间
	timeList2   []int     //第二轮响应时间
	timeList3   []int     //第三轮响应时间
	timeList    []int     //响应时间
	channel           = make(chan int64)
	done              = make(chan struct{})
	lock              = sync.Mutex{}

	iNum      int64 = 0
	firstNum  int64 //第一轮请求数
	secondNum int64 //第二轮请求数
	thirdNum  int64 //第三轮请求数
	one       int64 //第一轮并发数
	two       int64 //第二轮并发数
	three     int64 //第三轮并发数
	first     int64 //第一轮压测时长
	second    int64 //第二轮压测时长
	third     int64 //第三轮压测时长
	apiNum    int64
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
func printTime(useTime int64, count int64, iNum int64, okNum int64, timeList []int) {
	fmt.Println(len(timeList))
	fmt.Println("并发数: ", count)
	fmt.Println("请求数: ", iNum+count)
	fmt.Println("成功的数量：", len(timeList))
	fmt.Printf("\033[31m失败的数量：%v \033[0m \n", int(iNum+count)-len(timeList))
	fmt.Printf("\033[31m失败率：%.2f%v \033[0m \n", float64(int(iNum+count)-len(timeList))/float64(iNum+count)*100, "%")
	fmt.Println(fmt.Sprintf("运行耗时：%v秒", float64(useTime)/1000))
	fmt.Println("50%用户响应时间：" + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(fiftyRespTime(timeList)), float64(fiftyRespTime(timeList))/1e6))
	fmt.Println("90%用户响应时间：" + fmt.Sprintf("%.3f微秒 ≈ %v秒", float64(ninetyRespTime(timeList)), float64(ninetyRespTime(timeList))/1e6))
	fmt.Printf("最大响应时间：%.3f秒 \n", float64(maxRespTime(timeList))/1e6)
	fmt.Printf("最小响应时间：%.3f微秒 ≈ %v秒 \n", float64(minRespTime(timeList)), float64(minRespTime(timeList))/1e6)
	fmt.Printf("平均响应时间是:%.3f秒 \n", float64(sumRespTime(timeList))/float64(len(timeList))/1e6)
	fmt.Printf("QPS：%.3f \n", float64(count)/(float64(sumRespTime(timeList))/float64(len(timeList))/1e6))
}

func main() {
	var apiList = []string{"(1)BetSlipOdds", "(2)BetHistory", "(3)Balance", "(4)Bet", "(5)TransactionHistory", "(6)BetSlipOdds", "(7)BetLimit", "(8)SearchDefault", "(9)Search"}
	fmt.Println("接口列表：", apiList)
	fmt.Println("请输入要请求的接口编号如：1")
	_, _ = fmt.Scan(&apiNum)
	//var count int64
	//fmt.Println("请输入并发数量：")
	//_, _ = fmt.Scan(&count)
	//var one int64
	fmt.Println("第一轮压测用户数:")
	_, _ = fmt.Scan(&one)
	//var first int64
	fmt.Println("第一轮运行时长(秒):")
	_, _ = fmt.Scan(&first)
	//var two int64
	fmt.Println("第二轮压测用户数：")
	_, _ = fmt.Scan(&two)
	//var second int64
	fmt.Println("第二轮运行时长(秒):")
	_, _ = fmt.Scan(&second)
	//var three int64
	fmt.Println("第三轮压测用户数：")
	_, _ = fmt.Scan(&three)
	//var third int64
	fmt.Println("第三轮运行时长(秒):")
	_, _ = fmt.Scan(&third)
	startTime := time.Now().UnixNano() / 1e6
	fmt.Printf("开始时间：%v \n", startTime)
	do()
	endTime := time.Now().UnixNano() / 1e6
	fmt.Printf("结束时间：%v \n", endTime)
	fmt.Printf("总耗时：%.3f 秒 \n", float64(endTime-startTime)/1000-(10+10))
	//fmt.Println("总并发数:", one+two+three)
	//fmt.Println("总请求数: ", iNum)
	//fmt.Println("总成功的数量: ", firstOkNum+secondOkNum+thirdOkNum)
	//fmt.Printf("\033[31m总失败的数量: %v \033[0m \n", iNum-firstOkNum-secondOkNum-thirdOkNum)
	//fmt.Printf("\033[31m总失败率：%.2f%v \033[0m \n", float64(iNum-firstOkNum-secondOkNum-thirdOkNum)/float64(iNum)*100, "%")
	//runtime.GC()
	//_, _ = fmt.Scanf("h")
}

//, one int64, first int64, two int64, second int64, three int64, third int64
func do() {
	fmt.Println("第一轮压测开始...")
	fmt.Println(fmt.Sprintf("***首轮并发协程为%v***", one))
	wg.Add(int(one))
	for i := 0; i < int(one); i++ {
		switch apiNum {
		case 1:
			go Bet(first)
		case 2:
			go BetHistory(first)
		case 3:
			go Balance(first)
		case 4:
			go Bet(first)
		case 5:
			go TransactionHistory(first)
		case 6:
			go BetSlipOdds(first)
		case 7:
			go BetLimit(first)
		case 8:
			go SearchDefault(first)
		case 9:
			go Search(first)
		}
	}
	wg.Wait()
	firstNum = iNum
	firstOkNum = okNum
	timeList1 = timeList
	useTime1 = eT - sT
	fmt.Println("第二轮压测开始...")
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", two))
	<-time.After(10 * time.Second)
	wg.Add(int(two))
	for i := 0; i < int(two); i++ {
		switch apiNum {
		case 1:
			go MatchAndMarket(second)
		case 2:
			go BetHistory(second)
		case 3:
			go Balance(second)
		case 4:
			go Bet(second)
		case 5:
			go TransactionHistory(second)
		case 6:
			go BetSlipOdds(second)
		case 7:
			go BetLimit(second)
		case 8:
			go SearchDefault(second)
		case 9:
			go Search(second)
		}
	}
	wg.Wait()
	secondNum = iNum - firstNum
	secondOkNum = okNum - firstOkNum
	timeList2 = timeList[len(timeList1):]
	useTime2 = eT - sT
	fmt.Println("第三轮压测开始...")
	fmt.Println(fmt.Sprintf("***10秒后加压至%v协程***", three))
	<-time.After(10 * time.Second)
	wg.Add(int(three))
	for i := 0; i < int(three); i++ {
		switch apiNum {
		case 1:
			go MatchAndMarket(third)
		case 2:
			go BetHistory(third)
		case 3:
			go Balance(third)
		case 4:
			go Bet(third)
		case 5:
			go TransactionHistory(third)
		case 6:
			go BetSlipOdds(third)
		case 7:
			go BetLimit(third)
		case 8:
			go SearchDefault(third)
		case 9:
			go Search(third)
		}
	}
	wg.Wait()
	thirdNum = iNum - firstNum - secondNum
	thirdOkNum = okNum - firstOkNum - secondOkNum
	timeList3 = timeList[len(timeList1)+len(timeList2):]
	useTime3 = eT - sT
	fmt.Println("\033[33m***第一轮压测结果***\033[0m")
	printTime(useTime1, one, firstNum, firstOkNum, timeList1)
	fmt.Println("\033[35m----------------------- \033[0m")
	fmt.Println("\033[33m***第二轮压测结果***\033[0m")
	printTime(useTime2, two, secondNum, secondOkNum, timeList2)
	fmt.Println("\033[35m----------------------- \033[0m")
	fmt.Println("\033[33m***第三轮压测结果***\033[0m")
	printTime(useTime3, three, thirdNum, thirdOkNum, timeList3)
	fmt.Println("\033[35m----------------------- \033[0m")

}
func MatchAndMarket(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		filterRequest := &protouse.FilterRequest{Pager: &protouse.Pager{Page: 1, PageSize: 10},
			IsLive: 1, MarketTypes: []uint32{}, SportIds: []uint32{1},
			Times: []*timestamppb.Timestamp{}, MarketGroupType: 0,
			MatchIds: []string{}, OutcomeIds: []uint64{}, IsOutright: 0, Tournaments: []string{},
			Seasons: []string{}, CategoryIds: []uint32{}}
		byteS, _ := proto.Marshal(filterRequest)
		reqUrl := ""
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["Content-Type"] = "application/x-www-form-urlencoded, application/protobuf;proto=feedApiProto.FilterRequest"
		headers["Authorization"] = token
		headers["lang"] = "zh"
		headers["origin"] = "https://sports.aisport.live"
		res, err2 := req.SetHeaders(headers).Post(reqUrl, byteS)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		log.Printf("响应时间：%v毫秒", resTime)
		lock.Unlock()
		//channel <- resTime
		body, _ := res.Body()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var matchResponse protouse.MatchAndMarketResponse
		err1 := proto.Unmarshal(body, &matchResponse)
		if err1 != nil {
			log.Println("解析返回结构体异常: ", err1)
		}
		reBt, _ := json.Marshal(matchResponse)
		if res.StatusCode() == 200 && string(reBt) != "" {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常: ", string(reBt))
		}
		fmt.Println(string(reBt))
		res.Close()
		//log.Println("响应码：", res.StatusCode())
		//log.Println("接口返回：", matchResponse.MarketExtByTypeId)
	}

}
func BetHistory(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		rand.Seed(time.Now().UnixNano())
		daysLit := []int{-1, -7}
		rands := rand.Intn(2)
		betHistoryReq := &protouse.BetHistoryRequest{
			StartTime: timestamppb.New(time.Now().AddDate(0, 0, daysLit[rands])),
			EndTime:   timestamppb.New(time.Now()),
			Settled:   true,
			Pager: &protouse.Pager{
				TotalRecords: 1,
				Page:         1,
				PageSize:     10,
			},
		}
		//fmt.Println(betHistoryReq)
		byteS, _ := proto.Marshal(betHistoryReq)
		reqUrl := "http:/"
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["lang"] = "zh"
		headers["origin"] = "http://sports.aisport.live"
		headers["Content-Type"] = "application/x-www-form-urlencoded, application/protobuf;protos=sportBookProto.BetHistoryRequest"
		headers["Authorization"] = token
		res, err2 := req.SetHeaders(headers).Post(reqUrl, byteS)
		if err2 != nil {
			log.Println("请求异常：", err2)
		}
		log.Println("响应码：", res.StatusCode())
		defer res.Close()
		resTime, err3 := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		if err3 != nil {
			log.Println("转换响应时间异常：", err2)
		}
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		log.Printf("响应时间：%v毫秒\n", resTime)
		//channel <- resTime
		body, err4 := res.Body()
		if err4 != nil {
			log.Println("获取响应体异常：", err4)
		}
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var playOrders protouse.PlayerOrders
		err := proto.Unmarshal(body, &playOrders)
		if err != nil {
			log.Println("解析返回结构体异常: ", err)
		}
		reBt, _ := json.Marshal(playOrders)
		if res.StatusCode() == 200 && string(reBt) != "" {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常: ", string(reBt))
		}
		//byS, _ := json.Marshal(playOrders)
		log.Println("接口返回：", playOrders.Orders[0])

	}

}

func Balance(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for true {
		reqUrl := "http:/"
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["Authorization"] = token
		res, _ := req.SetHeaders(headers).Get(reqUrl)
		defer res.Close()
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
		atomic.AddInt64(&iNum, 1)
		var balance protouse.BalanceResponse
		err := proto.Unmarshal(body, &balance)
		if err != nil {
			log.Println("解析返回结构体异常: ", err)
		}
		reBt, _ := json.Marshal(balance)
		if res.StatusCode() == 200 && string(reBt) != "" {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常: ", string(reBt))
		}
		byS, _ := json.Marshal(balance)
		log.Println("接口返回：", string(byS))
		log.Println("响应码：", res.StatusCode())
	}

}
func Bet(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		var placeBetReq = &protouse.PlaceBetRequest{
			AcceptOddsChange: true,
			Selections: []*protouse.SelectionList{

				{
					MarketId:  "15396282",
					OutcomeId: "15396284",
					Odds:      "1.49",
				},
			},
			BetDetails: []*protouse.MultiLineDetail{
				{
					Type:  1,
					Stake: 7,
				},
			},
			OddsType: 0,
		}
		//fmt.Println(placeBetReq)
		byteS, _ := proto.Marshal(placeBetReq)
		reqUrl := "https://"
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["Content-Type"] = "application/x-www-form-urlencoded, application/protobuf;proto=sportBookProto.PlaceBetRequest"
		headers["Authorization"] = token
		//headers["origin"] = "https://sports.aisport.live"
		//headers["lang"] = "zh"
		res, _ := req.SetHeaders(headers).Post(reqUrl, byteS)
		defer res.Close()
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		log.Printf("响应时间：%v毫秒", resTime)
		//channel <- resTime
		body, _ := res.Body()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var betResponse protouse.Response
		err := proto.Unmarshal(body, &betResponse)
		if err != nil {
			log.Println("解析返回结构体异常: ", err)
		}
		reBt, _ := json.Marshal(betResponse)
		if res.StatusCode() == 200 && string(body) == "" {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常: ", string(reBt))
		}
		byS, _ := json.Marshal(betResponse)
		log.Println("接口返回：", string(byS))
		log.Println("响应码：", res.StatusCode())
	}
}

//钱包记录
func TransactionHistory(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		var transactionHistoryReq = &protouse.TransactionHistoryRequest{
			StartTime: timestamppb.New(time.Now().AddDate(0, -1, 0)),
			EndTime:   timestamppb.New(time.Now()),
			Pager: &protouse.Pager{
				TotalRecords: 1,
				Page:         1,
				PageSize:     10,
			},
		}
		byteS, _ := proto.Marshal(transactionHistoryReq)
		reqUrl := "http:/"
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["Content-Type"] = "application/x-www-form-urlencoded, application/protobuf;protos=sportBookProto.TransactionHistoryRequest"
		headers["Authorization"] = token
		res, _ := req.SetHeaders(headers).Post(reqUrl, byteS)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		defer res.Close()
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		//channel <- resTime
		body, _ := res.Body()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var transactionResponse protouse.Transaction
		err1 := proto.Unmarshal(body, &transactionResponse)
		if err1 != nil {
			log.Println("解析返回结构体异常: ", err1)
		}
		reBt, _ := json.Marshal(transactionResponse)
		if res.StatusCode() == 200 && string(reBt) != "" {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常: ", string(reBt))
		}
		byS, _ := json.Marshal(transactionResponse)
		log.Println("接口返回：", string(byS))
		log.Println("响应码：", res.StatusCode())
	}
}

func BetSlipOdds(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		var betSlipReq = &protouse.BetSlipRefreshRequest{
			OutcomeIds: []string{
				"11482241",
			},
			MarketIds: []string{
				"11482173",
			},
			Pager: &protouse.Pager{
				TotalRecords: 1,
				Page:         1,
				PageSize:     10,
			},
		}
		byteS, _ := proto.Marshal(betSlipReq)
		//data := make(map[string]interface{})
		//data["OutcomeIds"] = "9997088"
		//data["MarketIds"] = "9997079"
		reqUrl := "http:/"
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["Content-Type"] = "application/x-www-form-urlencoded, application/protobuf;protos=sportBookProto.BetSlipRefreshRequest"
		headers["Authorization"] = token
		res, _ := req.SetHeaders(headers).Post(reqUrl, byteS)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", resTime)
		defer res.Close()
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		//channel <- resTime
		body, _ := res.Body()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var betSlipResp protouse.BetSlipRefreshResponse
		err1 := proto.Unmarshal(body, &betSlipResp)
		if err1 != nil {
			log.Println("解析返回结构体异常: ", err1)
		}
		reBt, _ := json.Marshal(betSlipResp)
		if res.StatusCode() == 200 && string(reBt) != "" {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常: ", string(reBt))
		}
		byS, _ := json.Marshal(betSlipResp)
		log.Println("接口返回：", string(byS))
		log.Println("响应码：", res.StatusCode())
	}
}
func BetLimit(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		var betSlipReq = &protouse.PlaceBetRequest{
			AcceptOddsChange: true,
			Selections: []*protouse.SelectionList{

				{
					MarketId:  "11482173",
					OutcomeId: "11482241",
				},
			},
			BetDetails: []*protouse.MultiLineDetail{
				{
					Type: 1,
				},
			},
			OddsType: 0,
		}
		byteS, _ := proto.Marshal(betSlipReq)
		//data := make(map[string]interface{})
		//data["OutcomeIds"] = "9997088"
		//data["MarketIds"] = "9997079"
		reqUrl := "http"
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["Content-Type"] = "application/x-www-form-urlencoded, application/protobuf;protos=sportBookProto.PlaceBetRequest"
		headers["Authorization"] = token
		res, _ := req.SetHeaders(headers).Post(reqUrl, byteS)
		resTime, _ := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		log.Printf("响应时间：%v毫秒", resTime)
		defer res.Close()
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		//channel <- resTime
		body, _ := res.Body()
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var betSlipResp protouse.BetLimitResponse
		err1 := proto.Unmarshal(body, &betSlipResp)
		if err1 != nil {
			log.Println("解析返回结构体异常: ", err1)
		}
		reBt, _ := json.Marshal(betSlipResp)
		if res.StatusCode() == 200 && string(reBt) != "" {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常: ", string(reBt))
		}
		byS, _ := json.Marshal(betSlipResp)
		log.Println("接口返回：", string(byS))
		log.Println("响应码：", res.StatusCode())
	}
}

func SearchDefault(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		reqUrl := "http"
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["Authorization"] = token
		res, err2 := req.SetHeaders(headers).Get(reqUrl)
		if err2 != nil {
			log.Println("返回异常: ", err2)
		}
		log.Println("响应码：", res.StatusCode())
		resTime, err3 := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		if err3 != nil {
			fmt.Println("响应时间转换异常: ", err3)
		}
		log.Printf("响应时间：%v毫秒", resTime)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		//channel <- resTime
		body, err4 := res.Body()
		if err4 != nil {
			log.Println("解析响应消息体异常：", err4)
		}
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var betSlipResp protouse.SearchAndSportTypeResp
		err1 := proto.Unmarshal(body, &betSlipResp)
		if err1 != nil {
			log.Println("解析返回结构体异常: ", err1)
		}
		reBt, err5 := json.Marshal(betSlipResp)
		if err5 != nil {
			log.Println("将响应内容转换json格式时异常:", err5)
		}
		if res.StatusCode() == 200 && string(reBt) != "" {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常: ", string(reBt))
		}
		//byS, _ := json.Marshal(betSlipResp)
		log.Println("接口返回：", betSlipResp.Res[0])
		res.Close()
	}
}

func Search(times int64) {
	defer wg.Done()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("捕获的异常：", err)
		}
	}()
	sT = time.Now().UnixNano() / 1e6
	for {
		reqData := protouse.SearchSuggestionReq{
			Keyword: "足球",
		}
		bs, _ := proto.Marshal(&reqData)
		reqUrl := "http:"
		req := HttpRequest.NewRequest()
		headers := make(map[string]string)
		headers["Authorization"] = token
		headers["lang"] = "zh"
		headers["content-type"] = "application/x-www-form-urlencoded, application/protobuf;protos=localization.SearchSuggestionReq"
		res, err2 := req.SetHeaders(headers).Post(reqUrl, bs)
		if err2 != nil {
			log.Println("返回异常: ", err2)
		}
		log.Println("响应码：", res.StatusCode())
		resTime, err3 := strconv.ParseInt(strings.Split(res.Time(), "m")[0], 10, 64)
		if err3 != nil {
			fmt.Println("响应时间转换异常: ", err3)
		}
		log.Printf("响应时间：%v毫秒", resTime)
		lock.Lock()
		timeList = append(timeList, int(resTime*1000))
		lock.Unlock()
		//channel <- resTime
		body, err4 := res.Body()
		if err4 != nil {
			log.Println("解析响应消息体异常：", err4)
		}
		eT = time.Now().UnixNano() / 1e6
		if eT-sT > times*1000 {
			break
		}
		atomic.AddInt64(&iNum, 1)
		var betSlipResp protouse.SearchAndSportTypeResp
		err1 := proto.Unmarshal(body, &betSlipResp)
		if err1 != nil {
			log.Println("解析返回结构体异常: ", err1)
		}
		reBt, err5 := json.Marshal(betSlipResp)
		if err5 != nil {
			log.Println("将响应内容转换json格式时异常:", err5)
		}
		if res.StatusCode() == 200 && string(reBt) != "" {
			atomic.AddInt64(&okNum, 1)
		} else {
			log.Println("响应异常: ", string(reBt))
		}
		//byS, _ := json.Marshal(betSlipResp)
		log.Println("接口返回：", betSlipResp.Res[0])
		res.Close()
	}
}
