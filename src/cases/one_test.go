package main

import (
	"fmt"
	"gostudy/src/utils"
	"os"
	"testing"
)

func Test_One(t *testing.T) {
	path, _ := os.Getwd()
	fmt.Println(path)
	fmt.Println(utils.GetTestData("login.json", 0))

}
