package main

import (
	//strm "GoCode/work_script/mq"
	//request "GoCode/work_script/protos/proto"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"gostudy/src/protouse"
	"gostudy/src/protouse/protos"
	"gostudy/src/utils"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

var (
	resultAction = [6][2]int{{0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {1, 9}}
	MatchId      = make([]string, 200)
	wg           sync.WaitGroup
	nc           *nats.Conn
	rdb          *redis.Client
	js           = utils.JetStreamContext(nc)
	client       = &http.Client{}
	bunch        = []int{3, 4, 5, 6}
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       3,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Panicf("redis connect error: %v \n", err)
		return
	}
}

type Outcomes struct {
	Odds        string `json:"odds"`
	OutcomeId   string `json:"outcome_id"`
	SelectionId string `json:"selection_id"`
}

type Markets struct {
	MarketId   string     `json:"market_id"`
	MarketType string     `json:"market_type"`
	Outcomes   []Outcomes `json:"outcomes"`
}

type works struct {
	in   chan int
	done chan bool
}

type marketDetail struct {
	SportId    string    `json:"sport_id"`
	Category   string    `json:"category"`
	Tournament string    `json:"tournament"`
	Markets    []Markets `json:"markets"`
}

type Series struct {
	marketId    string
	marketType  string
	category    string
	tournament  string
	selectionId string
	matchId     string
	sportId     string
	OutcomeId   string
	Odds        string
}

type settleDetail struct {
	Odds        string  `json:"odds"`
	Outcomes    string  `json:"outcomes"`
	SettleState []int   `json:"settle_state"`
	Stake       float64 `json:"stake"`
}

type settleVerify struct {
	markets     []string
	selections  []string
	marketTypes []string
	matches     []string
	categories  []string
	sd          []settleDetail
}

type Result struct {
	Information  []settleDetail `json:"information"`
	To           float64        `json:"to"` // 总赔率
	Stake        float64        `json:"stake"`
	EstReturn    float64        `json:"est_return"` // 返回金额
	ExpectResult int            `json:"expect_result"`
	OrderResult  int            `json:"order_result"`
}

var requestData = struct {
	url           string
	method        string
	contentType   string
	authorization string
}{
	"",
	"POST",
	"application/x-www-form-urlencoded, application/protobuf;proto=feedApiProto.FilterRequest",
	"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI2MjE0MDUsImlhdCI6MTY1MjI2MTQwNSwiaXNzIjoic3BvcnRib29rX2FwaSIsInN1YiI6InNwb3J0IiwiT3BlcmF0b3JJZCI6MjUsIlNpdGVJZCI6NjQsIlVzZXJJZCI6MTUxMTc1NDQ0OCwiQ3VycmVuY3lJZCI6MSwiQWNjb3VudCI6ImNhcmxjYyJ9.vr_wHv-VEG9f2YGrWZFdv8TN1barOewVNjR2GReq7ESQA0iUlHN6oPwobmdqQF3dN6-GFDRxHvKgpo67H1ZAYQ",
}

func mathCombination(n int, m int) int {
	return factorial(n) / (factorial(n-m) * factorial(m))
}

func factorial(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

func combinationResult(n int, m int) [][]int {
	if m < 1 || m > n {
		fmt.Println("Illegal argument. Param m must between 1 and len(nums).")
		return [][]int{}
	}
	//保存最终结果的数组，总数直接通过数学公式计算
	result := make([][]int, 0, mathCombination(n, m))
	//保存每一个组合的索引的数组，1表示选中，0表示未选中
	index := make([]int, n)
	for i := 0; i < n; i++ {
		if i < m {
			index[i] = 1
		} else {
			index[i] = 0
		}
	}
	//第一个结果
	result = addTo(result, index)
	for {
		find := false
		//每次循环将第一次出现的 1 0 改为 0 1，同时将左侧的1移动到最左侧
		for i := 0; i < n-1; i++ {
			if index[i] == 1 && index[i+1] == 0 {
				find = true
				index[i], index[i+1] = 0, 1
				if i > 1 {
					moveOneToLeft(index[:i])
				}
				result = addTo(result, index)
				break
			}
		}
		//本次循环没有找到 1 0 ，说明已经取到了最后一种情况
		if !find {
			break
		}
	}
	return result
}

func addTo(arr [][]int, ele []int) [][]int {
	newEle := make([]int, len(ele))
	copy(newEle, ele)
	arr = append(arr, newEle)
	return arr
}

func moveOneToLeft(leftNums []int) {
	//计算有几个1
	sum := 0
	for i := 0; i < len(leftNums); i++ {
		if leftNums[i] == 1 {
			sum++
		}
	}
	//将前sum个改为1，之后的改为0
	for i := 0; i < len(leftNums); i++ {
		if i < sum {
			leftNums[i] = 1
		} else {
			leftNums[i] = 0
		}
	}
}

func findNumsByIndex(nums []int, index [][]int) [][]int {
	if len(index) == 0 {
		return [][]int{}
	}
	result := make([][]int, len(index))
	for i, v := range index {
		line := make([]int, 0)
		for j, v2 := range v {
			if v2 == 1 {
				line = append(line, nums[j])
			}
		}
		result[i] = line
	}
	return result
}

func creatWork(id int) works {
	w := works{
		in:   make(chan int),
		done: make(chan bool),
	}
	go doJob(id, w.in, w.done)
	return w
}

func Decimal(value float64) float64 {
	if value == math.Trunc(value) {
		return value
	}
	if len(strings.Split(fmt.Sprintf("%v", value), ".")[1]) > 2 {
		return math.Trunc(value*1e2+1) * 1e-2
	} else {
		return math.Trunc(value*1e2) * 1e-2

	}
}

func Truncate(f float64, prec int) float64 {
	// 不四舍五入保留指定位数
	n := strconv.FormatFloat(f, 'f', -1, 64)
	if n == "" {
		return 0
	}
	if prec >= len(n) {
		res, _ := strconv.ParseFloat(n, 64)
		return res
	}
	newn := strings.Split(n, ".")
	if len(newn) < 2 || prec >= len(newn[1]) {
		res, _ := strconv.ParseFloat(n, 64)
		return res
	}
	data := newn[0] + "." + newn[1][:prec]
	res, _ := strconv.ParseFloat(data, 64)
	return res
}

func calculateCombination(loop int, num []settleDetail) [][]settleDetail {
	nums := make([]int, 0)
	for i := 0; i < loop; i++ {
		nums = append(nums, i)
	}
	tmp := make([][]settleDetail, 0)
	for h := 2; h <= loop; h++ { // loop 3
		index := combinationResult(loop, h)
		hrs := findNumsByIndex(nums, index)
		// [[0 1] [0 2] [1 2] [0 1 2]
		for _, hr := range hrs { // [0 1]
			// 根据组合好的索引取出对应的值
			ss := make([]settleDetail, 0)
			for _, com := range hr {
				ss = append(ss, num[com])
			}
			tmp = append(tmp, ss)
		}
	}

	return tmp
}

func calculateOdd(data []settleDetail) (float64, [][]int) {
	var to float64 = 1
	var ar [][]int
	for _, e := range data {
		ar = append(ar, e.SettleState)
		fod, _ := strconv.ParseFloat(e.Odds, 64)
		to *= fod
	}
	return to, ar
}

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func MarshalBinary(d interface{}) ([]byte, error) {
	return json.Marshal(d)
}

func isIn(s []int, num int) bool {
	fm := make(map[int]int)
	for i, v := range s {
		fm[v] = i
	}
	if _, ok := fm[num]; ok {
		return true
	}
	return false
}

func handlerExpect(verify settleVerify, el int) {
	// 结算完成把结果写到redis
	var result Result
	resultSlice := make([]Result, 0)
	markets, _ := MarshalBinary(verify.markets)
	selections, _ := MarshalBinary(verify.selections)
	marketTypes, _ := MarshalBinary(verify.marketTypes)
	matches, _ := MarshalBinary(verify.matches)
	categories, _ := MarshalBinary(verify.categories)
	data := make(map[string]interface{})
	data["markets"] = markets
	data["selections"] = selections
	data["marketTypes"] = marketTypes
	data["matches"] = matches
	data["categories"] = categories
	com := calculateCombination(el, verify.sd)
	howCom := len(com)
	data["current_type_name"] = fmt.Sprintf("当前串关: [%d]串[%d]\n", el, howCom)
	data["total_stake"] = 1.69 * float64(howCom)
	log.Printf("当前串关: [%d]串[%d]\n", el, howCom)
	for oi, t := range com { // 4-1 11
		totalOdd, ar := calculateOdd(t)
		result.Information = t
		result.To = Truncate(totalOdd, 2)
		// Stake: 1.69 ar   [[0 2] [0 2]]  [[0 2] [1 9]] [[0 2] [1 9]]  [[0 2] [0 2] [1 9]]
		status := make([]int, 0) // 每一组状态
		for _, outside := range ar {
			// [{3.20 30128747 [3 4]}  {4.60 33342444 [3 4]}]
			result.Stake = 1.69
			av := outside[0]
			rv := outside[1]
			switch {
			case av == 0 && rv == 2:
				//result win  相当于发送settlement消息,并且赛果是赢的.
				status = append(status, 2)
			case av == 0 && rv == 3:
				// result HalfWon 相当于发送settlement消息,并且赛果是半赢的.
				status = append(status, 3)
			case av == 0 && rv == 4:
				// result Lost  相当于发送settlement消息,并且赛果是输的.
				status = append(status, 4)
			case av == 0 && rv == 5:
				// result HalfLost  相当于发送settlement消息,并且赛果是半输的.
				status = append(status, 5)
			case av == 0 && rv == 6:
				// result Void 相当于发送settlement消息,并且赛果是走水的
				status = append(status, 6)
			case av == 1 && rv == 9:
				// Cancel Cancel cancel 取消这个时间段内投注的订单.无论订单是否结算 对已经结算或者未结算都可以取消
				// 针对未结算时Cancel返回本金  针对已结算的Cancel扣回派奖 返回本金
				status = append(status, 9)
			default:
				break
			}
		}

		// 最后根据每个结果判断预期
		if isIn(status, 4) {
			// 有一个输,算全输
			result.EstReturn = 0
			result.ExpectResult = 4
			result.OrderResult = 4
		} else {
			/*
				假设投了10的5串1.每个赔率假设是2,[5 2 6 5 6, 3]
				第一个5: 10/2=5
				第二个2: 5*2=10
				第三个6: 10*1=10
				第四个5: 10/2=5
				第五个6: 5*1=5
				第六个3 5+(5/2)*(2-1)
			*/
			var repeat = make([]int, 0)
			var rMoney float64 = 1.69
			for index, s := range status {
				if s == 9 {
					repeat = append(repeat, s)
				}
				v, _ := strconv.ParseFloat(t[index].Odds, 64)
				switch s {
				case 2:
					rMoney *= v
					//lastTime.odds = append(lastTime.odds, t[index].Odds)
				case 3:
					rMoney += (1.69 / v) * (v - 1)
				case 5:
					rMoney /= v
				case 6:
					rMoney *= 1
				case 9:
					rMoney *= 1
				default:
					log.Println("status状态不存在!")
					break
				}
			}
			result.EstReturn = Decimal(rMoney)
			if len(repeat) == len(status) && result.EstReturn == 1.69 {
				// 取消 全部为9 并且本金=派奖
				result.ExpectResult = 9
				result.OrderResult = 9
			} else if result.EstReturn > 1.69 {
				// 赢 派奖大于本金
				result.ExpectResult = 2
				result.OrderResult = 2
			} else if result.EstReturn == 1.69 {
				// 走水 同时本金=派奖
				result.ExpectResult = 6
				result.OrderResult = 6
			} else if isIn(status, 5) && result.EstReturn > 0 && result.EstReturn < 1.69 {
				// 有一个或者多个半输 同时本金>派奖 派奖要大于0
				result.ExpectResult = 5
				result.OrderResult = 5
			}
			log.Printf("第 [%d] 次, rMoney=%v\n", oi, Decimal(rMoney))
		}
		resultSlice = append(resultSlice, result)
	}
	r, err := json.Marshal(resultSlice)
	if err != nil {
		log.Panicf("redis Marshal error: %v \n", err)
		return
	}
	data["outcome_information"] = r
	// 一次性保存多个hash字段值
	key := randStr(8)
	err = rdb.HMSet(key, data).Err()
	if err != nil {
		log.Panicf("set score failed error: %v \n", err)
		return
	}
}

func publish(js nats.JetStreamContext, subj string, f func(market, outcome uint64, result, action int32) []byte,
	market, outcome uint64, result, action int32) {
	// 返回发布的消息
	ack, err := js.Publish(subj, f(market, outcome, result, action))
	if err != nil {
		log.Fatalf("sigle_settle publish error: %v \n", err)
	}
	log.Printf("结算完成: Sequence --> %+v\n", ack.Sequence)
}

func testMsg(market, outcome uint64, result, action int32) []byte {
	var s = make([]*protos.Settle, 1)
	s = []*protos.Settle{
		{
			Producer: 1,
			Scope:    []uint32{1, 3},
			Market:   market,
			Outcomes: []*protos.OutcomeResult{
				{
					Outcome: outcome,
					Result:  protos.OutcomeResultCode(result),
					Shared:  1.0,
					Scope:   []uint32{1, 3},
				},
			},
			// 只针对cancel 一分钟内
			From: timestamp.New(time.Now().Add(-time.Minute)),
			To:   timestamp.New(time.Now().Add(time.Minute)),
		},
	}
	var settleShell = protos.SettleShell{
		Action:    protos.SettleAction(action), // Action 0～3 result Cancel RollbackCancel RollbackSettle
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
	} else {
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

func doJob(id int, c chan int, done chan bool) {
	// 等待任务结束
	for d := range c {
		log.Printf("goroutine %d 从channel:%v 收到的数据:%v\n", id+2, c, d)
		done <- true
	}
}

func RequestFirst(bodyData []byte) []int {
	// 根据不同MarketGroupType请求获取数据
	body, err := httpRequest(bodyData, "https://sports.9et.uk/api/v4/match_and_market")
	response := &protouse.MatchAndMarketResponse{}
	err = proto.Unmarshal(body, response)
	if err != nil {
		log.Fatalf("RequestFirst error: %v \n", err)
		return nil
	}
	for _, matches := range response.Matches {
		MatchId = append(MatchId, matches.MatchId)
	}
	noNoneVal := removeDuplicateElement(MatchId)
	return noNoneVal
}

func RequestSecond(gIndex int, bodyData []byte) (response *protouse.MatchAndMarketResponse) {
	body, err := httpRequest(bodyData, "https://sports.9et.uk/api/v4/match_and_market")
	response = &protouse.MatchAndMarketResponse{}
	err = proto.Unmarshal(body, response)
	if err != nil {
		log.Printf("第%d 协程 RequestSecond Unmarshal has error: %v \n", gIndex, err)
		return nil
	}
	return response
}

func getMarketsDetail(index int, match []int, wg *sync.WaitGroup) {
	defer wg.Done()
	var cc = make(map[int]marketDetail)
	for _, val := range match { // 一次串关的赛事
		resp := responseHandle(index, strconv.Itoa(val))
		var t marketDetail
		for _, m := range resp.Matches {
			t.SportId = strconv.Itoa(int(m.SportId))
			t.Category = m.Category
			t.Tournament = m.Tournament
			for _, market := range m.Markets {
				m := &Markets{
					MarketId:   market.MarketId,
					MarketType: market.MarketType,
				}
				for _, oc := range market.Outcomes {
					o := &Outcomes{
						OutcomeId:   oc.OutcomeId,
						Odds:        oc.Odds,
						SelectionId: oc.SelectionId,
					}
					m.Outcomes = append(m.Outcomes, *o)
				}
				t.Markets = append(t.Markets, *m)
			}
			cc[val] = t
		}
	}

	// 组合串关投注  noNoneVal 三串四：Trixie(11) 四串11: Yankee(12) 五串26: Super Yankee(13) 六串57: Heinz(14)
	var tt int
	currentBunch := len(cc) // 几场赛事 3
	keys := make([]int, 0)
	markets := make([]int, 0)
	for k, v := range cc {
		keys = append(keys, k)
		markets = append(markets, len(v.Markets))
	}
	theMinMarket := min(markets)
	switch currentBunch {
	case 3:
		tt = 11
	case 4:
		tt = 12
	case 5:
		tt = 13
	case 6:
		tt = 14
	}
	log.Printf("current_type: [%d] theMin[%d], 有%d场赛事, 赛事列表: %v\n", tt, theMinMarket, currentBunch, keys)
	for i := 0; i < theMinMarket; i++ { // theMin 100 Markets
		for j := 0; j < 2; j++ { // 每个Market下的outcome
			place := make([]*protouse.SelectionList, 0)
			var series Series
			seriesPlace := make([]Series, 0)
			for ii := 0; ii < currentBunch; ii++ { // 3 赛事长度
				// 每次循环根据当前组合赛事场数后在拿到每场的一个投注项
				test := cc[keys[ii]].Markets
				series.sportId = cc[keys[ii]].SportId
				series.category = cc[keys[ii]].Category
				series.tournament = cc[keys[ii]].Tournament
				series.marketId = test[i].MarketId
				series.marketType = test[i].MarketType
				series.selectionId = test[i].Outcomes[j].SelectionId
				series.OutcomeId = test[i].Outcomes[j].OutcomeId
				series.Odds = test[i].Outcomes[j].Odds
				series.matchId = strconv.Itoa(keys[ii])
				/*
					当前赛事: 2788539200, 第0次的Outcomes信息: [[<nil> OutcomeId:"29478117"  Odds:"1.73" OutcomeId:"29478085"  Odds:"1.88"]]
					当前赛事: 2787029000, 第0次的Outcomes信息: [[<nil> OutcomeId:"29386416"  Odds:"1.80" OutcomeId:"29386352"  Odds:"1.80"]]
					当前赛事: 2922885000, 第0次的Outcomes信息: [[<nil> OutcomeId:"30871673"  Odds:"5.20" OutcomeId:"30871674"  Odds:"1.13"]]
				*/
				selection := &protouse.SelectionList{
					OutcomeId: series.OutcomeId,
					Odds:      series.Odds,
				}
				place = append(place, selection)
				seriesPlace = append(seriesPlace, series)
			}
			// 投注逻辑
			var BetRequest = protouse.PlaceBetRequest{
				AcceptOddsChange: true,
				Selections:       place,
				BetDetails: []*protouse.MultiLineDetail{
					{
						Type:  protouse.OrderType(tt),
						Stake: 1.69,
					},
				},
			}
			b, err := proto.Marshal(&BetRequest)
			if err != nil {
				log.Printf("第[%d]->bet proto Marshal has error: %v \n", i, err)
				return
			}
			rep, _ := httpRequest(b, "https://sports.9et.uk/api/v4/bet")
			betResponse := &protouse.Response{}
			err = proto.Unmarshal(rep, betResponse)
			if err != nil {
				log.Printf("第[%d]->bet bet proto Unmarshal has error: %v \n", i, err)
				return
			}
			if betResponse.Status != 0 {
				log.Printf("第[%d]->bet 投注失败 error:  不继续结算%v \n", i, betResponse.Status)
				continue
			}
			fmt.Printf("第[%d]=========>: 投注完成开始进行结算: [%v]\n", i, &betResponse.Status)

			// 结算逻辑
			/*
				place
				[OutcomeId:"31479043"  Odds:"1.75" OutcomeId:"28097385"  Odds:"1.59" OutcomeId:"31752164"  Odds:"1.01"]
			*/
			var verify settleVerify
			for i := 0; i < len(place); i++ {
				var state settleDetail
				rand.Seed(time.Now().Unix())
				ra := resultAction[rand.Intn(len(resultAction))]
				md, _ := strconv.Atoi(seriesPlace[i].marketId)
				od, _ := strconv.Atoi(place[i].OutcomeId)
				// 投注金额: 88.69 action + Result状态  预期: 盘口Result order下的bet status 派奖金额
				state.Odds = seriesPlace[i].Odds
				state.Outcomes = seriesPlace[i].OutcomeId
				state.Stake = 1.69
				state.SettleState = append(state.SettleState, ra[0], ra[1])
				verify.markets = append(verify.markets, seriesPlace[i].marketId)
				verify.marketTypes = append(verify.marketTypes, seriesPlace[i].marketType)
				verify.selections = append(verify.selections, seriesPlace[i].selectionId)
				verify.matches = append(verify.matches, seriesPlace[i].matchId)
				verify.categories = append(verify.categories, seriesPlace[i].category)
				verify.sd = append(verify.sd, state)
				// 组合逻辑
				time.Sleep(time.Second)
				publish(js, "Settlement.Result", testMsg,
					uint64(md), uint64(od), int32(ra[1]), int32(ra[0]))

			}
			// 投注并结算一组后写入redis
			handlerExpect(verify, len(place))
		}
	}
	log.Printf("go程-[%d] 完成任务 \n", index)
}

func min(l []int) (min int) {
	min = l[0]
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return
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

func assembleSecondData(index int, mid string) (response *protouse.MatchAndMarketResponse) {
	// 根据matchIds Get Match Detail Response
	var FilterRequest = &protouse.FilterRequest{
		MatchIds: []string{mid},
	}
	b, err := proto.Marshal(FilterRequest)
	if err != nil {
		log.Printf("第%d 协程 assembleSecondData proto Marshal has error: %v \n", index, err)
		return nil
	}
	response = RequestSecond(index, b)
	return response
}

func assembleFirstData(mid string) (response *protouse.MatchAndMarketResponse) {
	allMatchId := make([]int, 0)
	var works [5]works
	for i := 0; i < 5; i++ {
		works[i] = creatWork(i)
	}
	for i := 0; i < 5; i++ {
		works[i].in <- 100
		var FilterRequest = &protouse.FilterRequest{
			IsLive:          1,
			MatchIds:        []string{},
			MarketGroupType: protouse.MarketGroupType(i + 2),
			MarketTypes:     []uint32{},
			SportIds:        []uint32{1},
			Times:           []*timestamp.Timestamp{},
			Pager:           &protouse.Pager{Page: 1, PageSize: 40},
		}
		b, err := proto.Marshal(FilterRequest)
		if err != nil {
			log.Fatalf("assembleFirstData proto Marshal has error: %v \n", err)
			return nil
		}
		matches := RequestFirst(b)
		allMatchId = append(allMatchId, matches...)
	}
	log.Printf("======= 数据全部接收完成 ->  即将开启%d个线程 ======== \n", len(allMatchId))
	// 组合串关投注  noNoneVal 三串四：Trixie(11) 四串11: Yankee(12) 五串26: Super Yankee(13) 六串57: Heinz(14)
	rand.Seed(time.Now().Unix())
	start := 0
	end := 0
	for end < len(allMatchId)-6 {
		ra := bunch[rand.Intn(len(bunch))]
		end += ra
		wg.Add(1)
		go func(start, end int) {
			match := allMatchId[start:end]
			getMarketsDetail(start, match, &wg) // match 下
		}(start, end)
		start = end
		break
	}
	log.Println(" --- 主线程完成, 等待其他任务 --- ")
	wg.Wait()
	return response
}

func responseHandle(index int, MatchId string) (response *protouse.MatchAndMarketResponse) {
	// 获取Match List Response get matchIds
	if MatchId == "" {
		response = assembleFirstData(MatchId)
	} else {
		response = assembleSecondData(index, MatchId)
	}
	return response
}

func main() {
	defer nc.Close()
	var MatchId string
	data := responseHandle(0, MatchId)
	fmt.Println("main 请求结果: ", data)
}
