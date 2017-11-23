package main

import (
	"fmt"
	"os"

	"github.com/dankraw/ssh-aliases/command"
)

// Version contains the binary version when built with -ldflags "-X main.Version=<some version>"
var Version string

func main() {
	cli, err := command.NewCLI(Version, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while configuring CLI:\n%v\n", err.Error())
		os.Exit(1)
	}
	err = cli.ApplyArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred during command execution:\n%v\n", err.Error())
		os.Exit(1)
	}
}
