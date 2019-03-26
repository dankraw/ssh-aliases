package command

import (
	"fmt"

	"io"

	"github.com/dankraw/ssh-aliases/compiler"
	"github.com/dankraw/ssh-aliases/config"
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

func (e *listCommand) execute(dir string, hosts []string) error {
	ctx, err := e.configReader.ReadConfigs(dir)
	if err != nil {
		return err
	}
	j := 0
	for _, s := range ctx.Sources {
		if len(s.Hosts) < 1 {
			continue
		}
		fileDelimiter := ""
		if j > 0 {
			fileDelimiter = "\n"
		}
		j++
		_, err = fmt.Fprint(e.writer, fileDelimiter+s.SourceName)
		if err != nil {
			return err
		}
		fmt.Fprintf(e.writer, " (%d):\n", len(s.Hosts))
		for _, h := range s.Hosts {
			results, err := e.compileHost(h, hosts)
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(e.writer, "\n "+h.AliasName)
			if err != nil {
				return err
			}
			fmt.Fprintf(e.writer, " (%d):\n", len(results))
			for _, r := range results {
				if r.HostName != "" {
					fmt.Fprintf(e.writer, "  %v: %v\n", r.Host, r.HostName)
				} else {
					fmt.Fprintf(e.writer, "  %v\n", r.Host)
				}
			}
		}
	}
	return nil
}

func (e *listCommand) compileHost(host compiler.ExpandingHostConfig, hosts []string) ([]compiler.HostEntity, error) {
	if host.IsRegexpHostDefinition() {
		return e.compiler.CompileRegexp(host, hosts)
	}
	return e.compiler.Compile(host)
}
