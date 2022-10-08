package utils

import (
	"runtime"
	"strings"
)

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	fuc := strings.Split(f.Name(), "/selenium/")[1]
	return fuc
}
