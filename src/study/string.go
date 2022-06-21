package main

import (
	"fmt"
	"unsafe"
)

func main() {
	a := "sa就"
	fmt.Println(string([]rune(a)[0:3]))
	//获取占据字节大小
	fmt.Println(unsafe.Sizeof(a))

	//string转成int
	//int, err := strconv.Atoi(string)
	//
	//string转成int64
	//int64, err := strconv.ParseInt(string, 10, 64)
	//
	//string转成uint64
	//uint64, err := strconv.ParseUint(string, 10, 64)
	//
	//int转成string
	//string := strconv.Itoa(int)
	//
	//int64转成string
	//string := strconv.FormatInt(int64, 10)
	//
	//uint64转成string(10进制)
	//string := strconv.FormatUint(uint64, 10)
	//
	//uint64转成string(16进制)
	//string := strconv.FormatUint(uint64, 16)
}
