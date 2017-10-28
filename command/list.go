package command

import (
	"fmt"
	"github.com/dankraw/ssh-aliases/compiler"
	"github.com/dankraw/ssh-aliases/config"
)

type ListCommand struct {
	configReader  *config.Reader
	configScanner *config.Scanner
	compiler      *compiler.Compiler
}

func NewListCommand() *ListCommand {
	return &ListCommand{
		configReader:  config.NewReader(),
		configScanner: config.NewScanner(),
		compiler:      compiler.NewCompiler(),
	}
}

func (e *ListCommand) List(dir string) error {
	files, err := e.configScanner.ScanDirectory(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		config, err := e.configReader.ReadConfig(f)
		if err != nil {
			return err
		}
		inputs, err := config.ToHostConfigInputs()
		if err != nil {
			return err
		}
		fmt.Printf("%v (definitions=%d):\n", f, len(inputs))
		for _, input := range inputs {
			results, err := e.compiler.Compile(input)
			fmt.Printf(" %v (compiled=%d):\n", input.AliasName, len(results))
			if err != nil {
				return err
			}
			for _, r := range results {
				fmt.Printf("  %v: %v\n", r.Host, r.HostName)
			}
		}
	}
	return nil
}
