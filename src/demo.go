package main

import (
	"fmt"
	"gostudy/src/utils"
)

const path = "C:"

func main() {
	fmt.Println(utils.GetTestData(path, "login.json", 0))

}
