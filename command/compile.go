package command

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

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
}

func NewCompileCommand(writer io.Writer) *CompileCommand {
	return &CompileCommand{
		indentation:  4,
		writer:       writer,
		configReader: config.NewReader(),
		compiler:     compiler.NewCompiler(),
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
	for _, input := range inputs {
		results, err := c.compiler.Compile(input)
		if err != nil {
			return err
		}
		for _, result := range results {
			c.printHostConfig(result)
		}
	}
	return nil
}

func (c *CompileCommand) printHostConfig(config HostConfigResult) {
	fmt.Fprintf(c.writer, "Host %v\n", config.Host)
	c.printHostConfigProperty("HostName", config.HostName)
	if config.HostConfig.IdentityFile != "" {
		c.printHostConfigProperty("IdentityFile", config.HostConfig.IdentityFile)
	}
	if config.HostConfig.Port != 0 {
		c.printHostConfigProperty("Port", config.HostConfig.Port)
	}
	fmt.Fprintln(c.writer)
}

func (c *CompileCommand) printHostConfigProperty(keyword string, value interface{}) {
	fmt.Fprintf(c.writer, "     %s %v\n", keyword, value)
}
