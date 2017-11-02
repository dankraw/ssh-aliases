package command

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"sort"

	"github.com/dankraw/ssh-aliases/compiler"
	"github.com/dankraw/ssh-aliases/config"
	. "github.com/dankraw/ssh-aliases/domain"
)

type CompileSaveCommand struct {
	file string
}

func NewCompileSaveCommand(file string) *CompileSaveCommand {
	return &CompileSaveCommand{
		file: file,
	}
}

func (c *CompileSaveCommand) Execute(dir string) error {
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
	sanitizer    *config.Sanitizer
}

func NewCompileCommand(writer io.Writer) *CompileCommand {
	return &CompileCommand{
		indentation:  4,
		writer:       writer,
		configReader: config.NewReader(),
		compiler:     compiler.NewCompiler(),
		validator:    compiler.NewValidator(),
		sanitizer:    config.NewSanitizer(),
	}
}

func (c *CompileCommand) Execute(dir string) error {
	config, err := c.configReader.ReadConfigs(dir)
	if err != nil {
		return err
	}
	inputs, err := config.ToHostConfigInputs()
	if err != nil {
		return err
	}
	allResults := []HostConfigResult{}
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

func (c *CompileCommand) printHostConfig(config HostConfigResult) {
	fmt.Fprintf(c.writer, "Host %v\n", config.Host)
	c.printHostConfigProperty("HostName", config.HostName)

	entries := config.HostConfig.ToHostConfigEntries()
	sort.Sort(ByHostConfigEntryKey(entries))

	for _, e := range entries {
		c.printHostConfigProperty(e.Key, e.Value)
	}
	fmt.Fprintln(c.writer)
}

func (c *CompileCommand) printHostConfigProperty(keyword string, value interface{}) {
	fmt.Fprintf(c.writer, "     %s %v\n", c.sanitizer.Sanitize(keyword), value)
}
