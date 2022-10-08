package main

import (
	"fmt"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	num := 10
	for i := 0; i < b.N; i++ {
		fmt.Println(num)
	}
}
