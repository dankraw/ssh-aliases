package command

import (
	"fmt"

	"io"

	"github.com/dankraw/ssh-aliases/compiler"
	"github.com/dankraw/ssh-aliases/config"
	"github.com/fatih/color"
)

type listCommand struct {
	writer        io.Writer
	configReader  *config.Reader
	configScanner *config.Scanner
	compiler      *compiler.Compiler
}

func newListCommand(writer io.Writer) *listCommand {
	return &listCommand{
		writer:        writer,
		configReader:  config.NewReader(),
		configScanner: config.NewScanner(),
		compiler:      compiler.NewCompiler(),
	}
}

func (e *listCommand) execute(dir string) error {
	files, err := e.configScanner.ScanDirectory(dir)
	if err != nil {
		return err
	}
	white := color.New(color.FgHiWhite)
	for i, f := range files {
		inputs, err := e.configReader.ReadConfig(f)
		if err != nil {
			return err
		}
		fileDelimiter := ""
		if i > 0 {
			fileDelimiter = "\n"
		}
		_, err = white.Fprint(e.writer, fileDelimiter+f)
		if err != nil {
			return err
		}
		fmt.Fprintf(e.writer, " (%d):\n", len(inputs))
		for _, input := range inputs {
			results, err := e.compiler.Compile(input)
			if err != nil {
				return err
			}
			_, err = white.Fprint(e.writer, "\n "+input.AliasName)
			if err != nil {
				return err
			}
			fmt.Fprintf(e.writer, " (%d):\n", len(results))
			for _, r := range results {
				fmt.Fprintf(e.writer, "  %v: %v\n", r.Host, r.HostName)
			}
		}
	}
	return nil
}
