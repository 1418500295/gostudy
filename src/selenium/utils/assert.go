package utils

import (
	"fmt"
	"strings"
)

var (
	OkNum        = 0
	FailNum      = 0
	FailCaseList []string
)

func AssertEqual(exp interface{}, actual interface{}, funcName string) {
	if exp == actual {
		fmt.Printf("\033[32m--- PASS   %v\n\033[0m", funcName)
		OkNum += 1
	} else {
		fmt.Printf("\033[31m--- FAIL   %v\n\033[0m", funcName)
		FailNum += 1
		FailCaseList = append(FailCaseList, funcName)
	}
}

func AssertContains(exp string, actual string, funcName string) {
	if strings.Contains(actual, exp) {
		fmt.Printf("\033[32m--- PASS   %v\n\033[0m", funcName)
		OkNum += 1
	} else {
		fmt.Printf("\033[31m--- FAIL   %v\n\033[0m", funcName)
		FailNum += 1
		FailCaseList = append(FailCaseList, funcName)
	}
}
