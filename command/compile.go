package command

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"os"

	"github.com/dankraw/ssh-aliases/compiler"
	"github.com/dankraw/ssh-aliases/config"
)

type CompileSaveCommand struct {
	file    string
	confirm *Confirm
}

func NewCompileSaveCommand(file string) *CompileSaveCommand {
	return &CompileSaveCommand{
		file:    file,
		confirm: NewConfirm(os.Stdin),
	}
}

func (c *CompileSaveCommand) Execute(dir string, force bool) error {
	if !force {
		confirmed, err := c.confirm.RequireConfirmationIfFileExists(c.file)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("File left unchanged.")
			return nil
		}
	}
	buffer := new(bytes.Buffer)
	err := NewCompileCommand(buffer).Execute(dir)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.file, buffer.Bytes(), 0600)
}

type CompileCommand struct {
	indentation  int
	writer       io.Writer
	configReader *config.Reader
	compiler     *compiler.Compiler
	validator    *compiler.Validator
}

func NewCompileCommand(writer io.Writer) *CompileCommand {
	return &CompileCommand{
		indentation:  4,
		writer:       writer,
		configReader: config.NewReader(),
		compiler:     compiler.NewCompiler(),
		validator:    compiler.NewValidator(),
	}
}

func (c *CompileCommand) Execute(dir string) error {
	config, err := c.configReader.ReadConfigs(dir)
	if err != nil {
		return err
	}
	inputs, err := config.ToExpandingHostConfigs()
	if err != nil {
		return err
	}
	allResults := []compiler.HostEntity{}
	for _, input := range inputs {
		results, err := c.compiler.Compile(input)
		if err != nil {
			return err
		}
		allResults = append(allResults, results...)
	}
	err = c.validator.ValidateResults(allResults)
	if err != nil {
		return err
	}
	for _, result := range allResults {
		c.printHostConfig(result)
	}
	return nil
}

func (c *CompileCommand) printHostConfig(config compiler.HostEntity) {
	fmt.Fprintf(c.writer, "Host %v\n", config.Host)
	c.printHostConfigProperty("HostName", config.HostName)

	for _, e := range config.Config {
		c.printHostConfigProperty(e.Key, e.Value)
	}
	fmt.Fprintln(c.writer)
}

func (c *CompileCommand) printHostConfigProperty(keyword string, value interface{}) {
	fmt.Fprintf(c.writer, "     %s %v\n", keyword, value)
}
