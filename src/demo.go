package main

import (
	"fmt"
	"runtime"
)

func main() {
	round := 5
	var wg sync.WaitGroup
	barrier := singleflight.Group{}
	wg.Add(round)
	for i := 0; i < round; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Printf("%d发起请求\n", i)
			// 启用10个协程模拟获取缓存操作
			value, err, _ := barrier.Do("get_rand_int", func() (interface{}, error) {
				fmt.Printf("%d正在运行\n", i)
				time.Sleep(2 * time.Second)
				fmt.Printf("%d执行完成\n", i)
				return send()
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%d协程返回结果是: %v", i, value)
			}
		}(i)
	}
	wg.Wait()
}



func main() {
	a()
	
	cy := cyclicbarrier.New(num)
	wg := sync.WaitGroup{}
	wg1 := sync.WaitGroup{}
	wg.Add(num)
	wg1.Add(1)
	for i := 0; i < num; i++ {
		go func(i int) {
			wg1.Wait()
			fmt.Println("准备开始。。。")
			time.Sleep(3 * time.Second)
			//fmt.Println(time.Now().UnixMilli())
			fmt.Printf("第%v个协程启动了\n", i)
			err := cy.Await(context.Background())
			if err != nil {
				return
			}
			wg.Done()
			fmt.Println("我在吃饭")
		}(i)
	}
	fmt.Println("开始")
	wg1.Done()
	wg.Wait()
	fmt.Println("结束")
}

func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func a() {
	fmt.Println(runFuncName())
}
