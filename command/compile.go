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

type compileSaveCommand struct {
	file    string
	confirm *confirm
}

func newCompileSaveCommand(file string) *compileSaveCommand {
	return &compileSaveCommand{
		file:    file,
		confirm: newConfirm(os.Stdin),
	}
}

func (c *compileSaveCommand) execute(dir string, force bool) error {
	if !force {
		confirmed, err := c.confirm.requireConfirmationIfFileExists(c.file)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("File left unchanged.")
			return nil
		}
	}
	buffer := new(bytes.Buffer)
	err := newCompileCommand(buffer).execute(dir)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.file, buffer.Bytes(), 0600)
}

type compileCommand struct {
	indentation  int
	writer       io.Writer
	configReader *config.Reader
	compiler     *compiler.Compiler
	validator    *compiler.Validator
}

func newCompileCommand(writer io.Writer) *compileCommand {
	return &compileCommand{
		indentation:  4,
		writer:       writer,
		configReader: config.NewReader(),
		compiler:     compiler.NewCompiler(),
		validator:    compiler.NewValidator(),
	}
}

func (c *compileCommand) execute(dir string) error {
	inputs, err := c.configReader.ReadConfigs(dir)
	if err != nil {
		return err
	}
	var allResults []compiler.HostEntity
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

func (c *compileCommand) printHostConfig(config compiler.HostEntity) {
	fmt.Fprintf(c.writer, "Host %v\n", config.Host)
	c.printHostConfigProperty("HostName", config.HostName)

	for _, e := range config.Config {
		c.printHostConfigProperty(e.Key, e.Value)
	}
	fmt.Fprintln(c.writer)
}

func (c *compileCommand) printHostConfigProperty(keyword string, value interface{}) {
	fmt.Fprintf(c.writer, "     %s %v\n", keyword, value)
}
