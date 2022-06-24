package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	//字符串拼接
	a := "sada"
	b := "1231"
	var buf bytes.Buffer
	buf.WriteString(a)
	buf.WriteString(b)
	s := buf.String()
	fmt.Println(s)

	var builder strings.Builder
	builder.WriteString(a)
	builder.WriteString(b)
	s1 := builder.String()
	fmt.Println(s1)

}
