package main

import (
	"fmt"
	"os/exec"
)

func main() {
	out, err := exec.Command("main", "2==2").Output()
	if err == nil {
		fmt.Println(string(out))
	}
}
