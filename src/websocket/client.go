package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"sync"
)

var origin = "http://127.0.0.1:8080/"
var url = "ws://127.0.0.1:8080/echo"
var wg = sync.WaitGroup{}

func run() {
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go Req()
	}
	wg.Wait()

}
func main() {
	run()
}

func Req() {
	defer wg.Done()
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte("hello, world!你好")
	_, err = ws.Write(message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Send: %s\n", message)

	var msg = make([]byte, 512)
	m, err := ws.Read(msg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Receive: %s\n", msg[:m])

	_ = ws.Close() //关闭连接

}
