package main

import (
	"fmt"
	"runtime"
)

func main() {
	a()
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
