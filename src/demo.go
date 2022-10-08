package main

import (
	"fmt"
	"os/exec"
)

var ()

func main() {
	cmd := exec.Command("bash", "-c", "cd ..")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))

	cmd1 := exec.Command("bash", "-c", "go test -json | go-test-report -o ../test_report.html")
	out1, err1 := cmd1.Output()
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println(string(out1))
}
