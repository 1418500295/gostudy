package main

import "os/exec"

func main() {
	exec.Command("cd ./cases/")
	exec.Command("go test -v")
}
