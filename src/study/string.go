package main

import "fmt"

func main() {
	a := "sa就是大概不会"
	fmt.Println(string([]rune(a)[0:3]))
}
