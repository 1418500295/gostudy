package main

import (
	"fmt"
	"gostudy/src/utils"
)

const path = "C:\\Users\\AA\\go\\src\\gostudy\\src"

func main() {
	fmt.Println(utils.GetTestData(path, "login.json", 0))

}
