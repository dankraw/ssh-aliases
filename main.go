package main

import (
	"fmt"
	"os"

	. "github.com/dankraw/ssh-aliases/command"
)

var VERSION string

func main() {
	err := NewCLI(VERSION).ConfigureCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred:\n%v\n", err.Error())
		os.Exit(1)
	}
}
