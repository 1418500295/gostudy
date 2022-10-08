package main

import (
	strm "GoCode/work_script/mq"
	request "GoCode/work_script/protos/proto"
	"bytes"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// MatchId resultAction = [6][2]int{{0,2}, {0,3}, {0,4}, {0,5}, {0,6}, {1,9}}
	MatchId  = make([]string, 100)
	wg sync.WaitGroup
	nc *nats.Conn
	rdb *redis.Client
	//js = strm.JetStreamContext(nc)
	client =  &http.Client{}
)

func init()  {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       2,  // use default DB
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Panicf("redis connect error: %v \n", err)
		return
	}
}

type Outcomes struct {
	Odds string `json:"odds"`
	OutcomeId string `json:"outcome_id"`
	SelectionId string `json:"selection_id"`
}

type Markets struct {
	MarketId string `json:"market_id"`
	MarketType string `json:"market_type"`
	Outcomes []Outcomes `json:"outcomes"`
}

var requestData = struct {
	url    string
	method string
	contentType string
	authorization string
}{
	"",
	"POST",
	"application/x-www-form-urlencoded, application/protobuf;proto=feedApiProto.FilterRequest",
	"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTcyNzAxOTUsImlhdCI6MTY1NjkxMDE5NSwiaXNzIjoic3BvcnRib29rX2FwaSIsInN1YiI6InNwb3J0IiwiT3BlcmF0b3JJZCI6MjUsIlNpdGVJZCI6NjQsIlVzZXJJZCI6MTUxMTc1NDQ0OCwiQ3VycmVuY3lJZCI6MSwiQWNjb3VudCI6ImNhcmxjYyJ9.IoHkt962pJ1RyOH4vCjX0FOAgeczBN3bSDneiw34BDbti91YNyM5ngshGhTvXZ5XQE9v2Ts5ASRvKtLxlzx0tg",
}

type marketDetail struct{
	SportId string `json:"sport_id"`
	Category string `json:"category"`
	Tournament string `json:"tournament"`
	Markets []Markets `json:"markets"`
}

func Decimal(value float64) float64 {
	if value == math.Trunc(value){
		return value
	}
	if len(strings.Split(fmt.Sprintf("%v", value), ".")[1]) > 2 {
		return math.Trunc(value*1e2+1) * 1e-2
	}else {
		return math.Trunc(value * 1e2) * 1e-2

	}
}

func handlerExpect(expect interface{})  {
	// 结算完成把结果写到redis
	data := make(map[string]interface{})
	val := reflect.ValueOf(expect)
	kd := val.Kind()
	if kd != reflect.Struct{
		log.Panicf("not a Struct\n")
		return
	}
	ac := val.FieldByName("Action")
	rt := val.FieldByName("Result")
	acVal := ac.Interface()
	rtVal := rt.Interface()
	av := acVal.(request.SettleAction)
	rv := rtVal.(request.OutcomeResultCode)
	marketInfo := val.FieldByName("MarketInfo")
	marketVal := marketInfo.Interface()
	mv := marketVal.(map[string]string)
	data["odd"] = mv["odd"]
	float,_ := strconv.ParseFloat(mv["odd"],64)
	switch  {
	// Stake: 20.69
	case av == 0 && rv == 2:
		//result win  相当于发送settlement消息,并且赛果是赢的.
		data["expect_result"] = 2
		data["est_return"] = Decimal(20.69 * float)
		data["order_result"] = 2
	case av == 0 && rv == 3:
		// result HalfWon 相当于发送settlement消息,并且赛果是半赢的.
		data["expect_result"] = 3
		data["est_return"] = Decimal(20.69 + (20.69/2) * (float-1)) //（1.69/2）*（float-1）
		data["order_result"] = 3
	case av == 0 && rv == 4:
		// result Lost  相当于发送settlement消息,并且赛果是输的.
		data["expect_result"] = 4
		data["est_return"] = 0
		data["order_result"] = 4
	case av == 0 && rv == 5:
		// result HalfLost  相当于发送settlement消息,并且赛果是半输的.
		data["expect_result"] = 5
		data["est_return"] = Decimal(20.69/2)
		data["order_result"] = 5
	case av == 0 && rv == 6:
		// result Void 相当于发送settlement消息,并且赛果是走水的
		data["expect_result"] = 6
		data["est_return"] = 20.69
		data["order_result"] = 6
	case av == 1 && rv == 9:
		// Cancel Cancel cancel 取消这个时间段内投注的订单.无论订单是否结算 对已经结算或者未结算都可以取消
		data["expect_result"] = 9
		data["est_return"] = 20.69   // 针对未结算时Cancel返回本金  针对已结算的Cancel扣回派奖 返回本金
		data["order_result"] = 9
	default:
		break
	}
	data["stake"] = 20.69
	data["action"] = ac.Int()
	data["reality_result"] = rt.Int()
	data["sport_id"] = mv["sport_id"]
	data["category"] = mv["category"]
	data["market_id"] = mv["market_id"]
	data["market_type"] = mv["market_type"]
	data["selection_id"] = mv["selection_id"]
	data["outcome_id"] = mv["outcome_id"]

	log.Printf("写入redis信息outcome: %v, 结算结果: %v-%v \n", data["outcome_id"], acVal, rtVal)
	// 一次性保存多个hash字段值
	err := rdb.HMSet(mv["outcome_id"], data).Err()
	if err != nil {
		log.Panicf("set score failed error: %v \n", err)
		return
	}
}

func publish(js nats.JetStreamContext,subj string, f func(market,outcome uint64, result, action int32) []byte,
	market, outcome uint64, result, action int32, expect interface{}) {
	// 返回发布的消息
	ack, err := js.Publish(subj, f(market, outcome,result, action))
	if err != nil {
		log.Fatalf("sigle_settle publish error: %v  ack: %v\n", err, ack)
	}
	log.Printf("当前结算outcome: %v, 结算结果: %v-%v \n", outcome, action, result)
	handlerExpect(expect)
	//time.Sleep(time.Second * 10)
}

func testMsg(market, outcome uint64, result, action int32) []byte {
	var s  = make([]*request.Settle, 1)
	s = []*request.Settle{
		{
			Producer: 1,
			Scope:[]uint32{1, 3},
			Market: market,
			Outcomes: []*request.OutcomeResult{
				{
					Outcome: outcome,
					Result: request.OutcomeResultCode(result),
					Shared:  1.0,
					Scope:   []uint32{1, 3},
				},
			},
			// 只针对cancel 一分钟内
			From: timestamp.New(time.Now().Add(-time.Minute)),
			To: timestamp.New(time.Now().Add(time.Minute)),
		},
	}
	var settleShell = request.SettleShell{
		Action: request.SettleAction(action), // Action 0～3 result Cancel RollbackCancel RollbackSettle
		Settle:    s,
		Timestamp: timestamp.New(time.Now().Add(time.Millisecond)),
	}
	b, err := proto.Marshal(&settleShell)
	if err != nil {
		log.Fatalf("sigle_settle Marshal error: %v \n", err)
		return nil
	}
	return b
}

func httpRequest(bodyData []byte, url string) ([]byte, error) {
	req, err := http.NewRequest(requestData.method, url, bytes.NewBuffer(bodyData))
	if err != nil {
		log.Fatalf("httpRequest->NewRequest 错误信息: %v \n", err)
		return nil, err
	}
	suffix := strings.HasSuffix(url, "bet")
	if suffix {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded, application/protobuf;proto=sportBookProto.PlaceBetRequest")
	}else {
		req.Header.Set("Content-Type", requestData.contentType)
	}
	req.Header.Set("authorization", requestData.authorization)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("httpRequest->client.Do 错误信息: %v \n", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func RequestFirst(bodyData []byte) (data string, err error) {
	body, err :=  httpRequest(bodyData, "https://sports.aisport.live/api/v4/match_and_market")
	response := &request.MatchAndMarketResponse{}
	err = proto.Unmarshal(body, response)
	if err != nil {
		log.Fatalf("RequestFirst error: %v \n", err)
		return "", err
	}
	for _, matches := range response.Matches{
		MatchId = append(MatchId, matches.MatchId)
	}
	// result, err := json.Marshal(response.MarketExtByTypeId)  转json字符串 string(result))
	// result, err := json.Marshal(response.Matches)
	noNoneVal := removeDuplicateElement(MatchId)
	log.Printf("即将开启%d个线程 \n", len(noNoneVal))
	wg.Add(len(noNoneVal))
	for index, val := range noNoneVal{
		if val != 0 {
			// 取出每个赛事下的MarketExtByTypeId  Matches->Markets->[MarketId,Odds,OutcomeId]
			go getMarketsDetail(index, strconv.Itoa(val), &wg)
			//getMarketsDetail(index, strconv.Itoa(val), &wg)
		}
	}
	wg.Wait()
	log.Println(" --- 主线程完成, 等待其他任务 --- ")
	return "RequestFirst finish", nil
}

func RequestSecond(gIndex int, bodyData []byte) (data string, err error) {
	body, err :=  httpRequest(bodyData, "https://sports.aisport.live/api/v4/match_and_market")
	response := &request.MatchAndMarketResponse{}
	err = proto.Unmarshal(body, response)
	if err != nil {
		log.Printf("第%d 协程 RequestSecond Unmarshal has error: %v \n", gIndex, err)
		return "", err
	}
	var t  marketDetail
	for _, m := range response.Matches{  // 1
		t.SportId = strconv.Itoa(int(m.SportId))
		t.Category = m.Category
		t.Tournament = m.Tournament
		for _, market := range m.Markets{  // 131
			m := &Markets{
				MarketId: market.MarketId,
				MarketType: market.MarketType,
			}
			for _, oc := range market.Outcomes{
				o := &Outcomes{
					OutcomeId: oc.OutcomeId,
					Odds: oc.Odds,
					SelectionId: oc.SelectionId,
				}
				m.Outcomes = append(m.Outcomes, *o)
			}
			t.Markets = append(t.Markets, *m)
		}

		// 单注投注数据构造
		for _, out := range t.Markets{
			for _, odd := range out.Outcomes{
				var BetRequest =  request.PlaceBetRequest{
					AcceptOddsChange: true,
					Selections: []*request.SelectionList{
						{
							OutcomeId: odd.OutcomeId,
							Odds: odd.Odds,
						},
					},
					BetDetails: []*request.MultiLineDetail{
						{
							Type: 1,  // 2串 2 3串 3
							Stake: 1.11,
						},
					},
				}
				fmt.Printf("%v,%v\n", odd.OutcomeId,odd.Odds)
				_, err := proto.Marshal(&BetRequest)
				if err != nil {
					log.Printf("第%d 协程 bet proto Marshal has error: %v \n", gIndex, err)
					return "", err
				}
				////log.Printf("第%d 协程 OutcomeId: %v , MarketId: %v, MarketType: %v, Category: %v, Tournament: %v",
				////	gIndex, odd.OutcomeId, out.MarketId, out.MarketType, t.Category, t.Tournament)
				//rep, _ := httpRequest(b, "https://sports.aisport.live/api/v4/bet")
				//betResponse := &request.Response{}
				//err = proto.Unmarshal(rep, betResponse)
				//if err != nil {
				//	log.Printf("第%d 协程 bet proto Unmarshal has error: %v \n", gIndex, err)
				//	return "", err
				//}
				//if betResponse.Status != 0 {
				//	log.Printf("第%d 协程 投注失败 error:  不继续结算%v \n", gIndex, betResponse.Status)
				//	continue
				//}
				//for _, val := range BetRequest.Selections{
				//	fmt.Printf("%v,%v\n",val.OutcomeId,val.Odds)
				//}
				//// 结算 2-8 action 0-3
				//rand.Seed(time.Now().Unix())
				//ra := resultAction[rand.Intn(len(resultAction))]
				//md, _ := strconv.Atoi(out.MarketId)
				//od, _ := strconv.Atoi(odd.OutcomeId)
				//// 投注金额: 88.69 action + Result状态  预期: 盘口Result order下的bet status 派奖金额
				//var expect = struct{
				//	MarketInfo map[string]string
				//	Action request.SettleAction
				//	Result request.OutcomeResultCode
				//}{
				//	Action: request.SettleAction(ra[0]),
				//	Result: request.OutcomeResultCode(ra[1]),
				//	MarketInfo: map[string]string{
				//		"sport_id": strconv.Itoa(int(m.SportId)),
				//		"category": m.Category,
				//		"market_id": out.MarketId,
				//		"odd": odd.Odds,
				//		"market_type":  out.MarketType,
				//		"selection_id": odd.SelectionId,
				//		"outcome_id": odd.OutcomeId,
				//	},
				//}
				publish(js, "Settlement.Result", testMsg,
				//	uint64(md), uint64(od), int32(ra[1]), int32(ra[0]), expect)
			}
		}
	}
	return fmt.Sprintf("RequestSecond 协程 %d finish", gIndex), nil
}

func betTask(r *request.PlaceBetRequest)  {
	b, err := proto.Marshal(r)
	if err != nil {
		fmt.Println("bet proto Marshal has error: ", err)
		return
	}
	rep, _ := httpRequest(b, "https://sports.9et.uk/api/v4/bet")
	betResponse := &request.Response{}
	err = proto.Unmarshal(rep, betResponse)
	if err != nil {
		fmt.Println("bet proto Unmarshal has error:", err)
		return
	}
	fmt.Println("bet success ", betResponse.Message, betResponse.Status)
	//result, _ := json.Marshal(t)
	//fmt.Println("marketDetail: ", string(result))
}

func getMarketsDetail(index int,MarketId string,  wg *sync.WaitGroup)  {
	defer wg.Done()
	// index go程
	resp := responseHandle(index, MarketId)
	log.Printf("第%d个go程 完成任务: %v \n", index, resp)
}

func removeDuplicateElement(languages []string) []int {
	result := make([]int, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if item == "" {
			continue
		}
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			val, _ := strconv.Atoi(item)
			result = append(result, val)
		}
	}
	return result
}

func sendSettle()  {
	var nc *nats.Conn
	js := strm.JetStreamContext(nc)
	defer nc.Close()
	publish(js, "Settlement.Result", testMsg, 27175152, 27175216, 2, 0, nil)
}

func assembleSecondData(index int,mid string) (repData string) {
	// 根据matchIds Get Match Detail Response
	var FilterRequest = &request.FilterRequest{
		MatchIds: []string{mid},
	}
	b, err := proto.Marshal(FilterRequest)
	if err != nil {
		log.Printf("第%d 协程 assembleSecondData proto Marshal has error: %v \n", index, err)
		return fmt.Sprintf("assembleSecondData proto Marshal has error: %v \n", err)
	}
	repData, _ = RequestSecond(index, b)
	return repData
}

func assembleFirstData(mid string) (repData string) {
	var FilterRequest = &request.FilterRequest{
		IsLive: 1,
		MatchIds: []string{},
		MarketGroupType: 2,
		MarketTypes: []uint32{},
		SportIds: []uint32{},  //  1 足球  MarketGroupType 1～6  2 篮球 MarketGroupType 11
		Times: []*timestamp.Timestamp{},
		Pager: &request.Pager{Page: 1, PageSize: 200},  // eg 500场赛事
	}
	b, err := proto.Marshal(FilterRequest)
	if err != nil {
		log.Fatalf("assembleFirstData proto Marshal has error: %v \n", err)
		return ""
	}
	repData, _ = RequestFirst(b)
	return repData
}

func responseHandle(index int,MatchId string) (rsp string) {
	// 获取Match List Response get matchIds
	if MatchId == "" {
		rsp = assembleFirstData(MatchId)
	}else {
		rsp = assembleSecondData(index, MatchId)
	}
	return rsp
}

func main() {
	defer nc.Close()
	var MatchId string
	data := responseHandle(0,MatchId)
	fmt.Println("main 请求结果: ", data)
}