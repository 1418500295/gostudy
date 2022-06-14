package main

import (
	"fmt"
	"sort"
)

func main() {
	a := []int{1, 34, 2}
	//升序
	sort.Ints(a)
	fmt.Println(a)
	//降序
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	fmt.Println(a)
}
