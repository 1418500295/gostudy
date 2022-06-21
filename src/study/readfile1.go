package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const path = "C:\\Users\\AA\\go\\src\\gostudy\\src\\testdata\\login.json"

func main() {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	r := bufio.NewReader(fi)
	var chunks []byte

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println(err)
		}
		if 0 == n {
			break
		}
		chunks = append(chunks, buf...)
		var jsonData []map[string]interface{}
		err1 := json.Unmarshal(chunks[:n], &jsonData)
		if err1 != nil {
			fmt.Println(err1)
		}
		fmt.Println(jsonData)
	}

	//fmt.Println(string(chunks))
}
func Read1() string {
	//获得一个file
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("read fail")
		return ""
	}

	//把file读取到缓冲区中
	defer f.Close()
	var chunk []byte
	buf := make([]byte, 1024)

	for {
		//从file读取到buf中
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read buf fail", err)
			return ""
		}
		//说明读取结束
		if n == 0 {
			break
		}
		//读取到最终的缓冲区中
		chunk = append(chunk, buf[:n]...)
	}

	return string(chunk)
	//fmt.Println(string(chunk))
}
