package command

import (
	"fmt"

	"io"

	"github.com/dankraw/ssh-aliases/compiler"
	"github.com/dankraw/ssh-aliases/config"
	"github.com/fatih/color"
)

type ListCommand struct {
	writer        io.Writer
	configReader  *config.Reader
	configScanner *config.Scanner
	compiler      *compiler.Compiler
}

func NewListCommand(writer io.Writer) *ListCommand {
	return &ListCommand{
		writer:        writer,
		configReader:  config.NewReader(),
		configScanner: config.NewScanner(),
		compiler:      compiler.NewCompiler(),
	}
}

func (e *ListCommand) Execute(dir string) error {
	files, err := e.configScanner.ScanDirectory(dir)
	if err != nil {
		return err
	}
	white := color.New(color.FgHiWhite)
	for i, f := range files {
		config, err := e.configReader.ReadConfig(f)
		if err != nil {
			return err
		}
		inputs, err := config.ToHostConfigInputs()
		if err != nil {
			return err
		}
		file_delimiter := ""
		if i > 0 {
			file_delimiter = "\n"
		}
		white.Fprint(e.writer, file_delimiter+f)
		fmt.Fprintf(e.writer, " (%d):\n", len(inputs))
		for _, input := range inputs {
			results, err := e.compiler.Compile(input)
			white.Fprint(e.writer, "\n "+input.AliasName)
			fmt.Fprintf(e.writer, " (%d):\n", len(results))
			if err != nil {
				return err
			}
			for _, r := range results {
				fmt.Fprintf(e.writer, "  %v: %v\n", r.Host, r.HostName)
			}
		}
	}
	return nil
}
