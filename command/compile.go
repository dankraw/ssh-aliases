package command

import (
	"bytes"
	"fmt"
	"io"
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

func (c *compileSaveCommand) execute(dir string, force bool, hosts []string) error {
	if !force {
		confirmed, err := c.confirm.requireConfirmationIfFileExists(c.file)
		if err != nil {
			return err
		}
		if confirmed {
			fmt.Printf("Writing changes to %s", c.file)
		} else {
			fmt.Printf("Exiting without writing changes to %s", c.file)
			return nil
		}
	}
	buffer := new(bytes.Buffer)
	err := newCompileCommand(buffer).execute(dir, hosts)
	if err != nil {
		return err
	}
	return os.WriteFile(c.file, buffer.Bytes(), 0o600)
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

func (c *compileCommand) execute(dir string, hosts []string) error {
	ctx, err := c.configReader.ReadConfigs(dir)
	if err != nil {
		return err
	}
	var allResults []compiler.HostEntity
	for _, s := range ctx.Sources {
		for _, h := range s.Hosts {
			results, err := c.compileHost(h, hosts)
			if err != nil {
				return err
			}
			allResults = append(allResults, results...)
		}
	}
	err = c.validator.ValidateResults(allResults)
	if err != nil {
		return err
	}
	for _, result := range allResults {
		err = c.printHostConfig(result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *compileCommand) compileHost(host compiler.ExpandingHostConfig, hosts []string) ([]compiler.HostEntity, error) {
	if host.IsRegexpHostDefinition() {
		result, err := c.compiler.CompileRegexp(host, hosts)
		return result, err
	}
	return c.compiler.Compile(host)
}

func (c *compileCommand) printHostConfig(cfg compiler.HostEntity) error {
	_, err := fmt.Fprintf(c.writer, "Host %v\n", cfg.Host)
	if err != nil {
		return err
	}
	if cfg.HostName != "" {
		err = c.printHostConfigProperty("HostName", cfg.HostName)
		if err != nil {
			return err
		}
	}
	for _, e := range cfg.Config {
		err = c.printHostConfigProperty(e.Key, e.Value)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(c.writer)
	return err
}

func (c *compileCommand) printHostConfigProperty(keyword string, value interface{}) error {
	_, err := fmt.Fprintf(c.writer, "     %s %v\n", keyword, value)
	return err
}
