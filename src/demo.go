package main

import (
	"fmt"
	"runtime"
)

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
