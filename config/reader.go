package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/dankraw/ssh-aliases/compiler"
)

// Reader is able to read directories and files and return inputs for ssh-aliases compiler
type Reader struct {
	decoder *decoder
	scanner *Scanner
}

// NewReader returns new instance of Reader
func NewReader() *Reader {
	return &Reader{
		decoder: newDecoder(),
		scanner: NewScanner(),
	}
}

// ReadConfigs processes the input directory and returns inputs for ssh-aliases compiler
func (e *Reader) ReadConfigs(dir string) (compiler.InputContext, error) {
	files, err := e.scanner.ScanDirectory(dir)
	if err != nil {
		return compiler.InputContext{}, err
	}
	var sources []rawContextSource
	for _, f := range files {
		c, err := e.decodeFile(f)
		if err != nil {
			return compiler.InputContext{}, errors.Wrap(err, fmt.Sprintf("failed parsing %s", f))
		}
		if len(c.Hosts) < 1 && len(c.RawConfigs) < 1 && len(c.Values) < 1 {
			continue
		}
		rawSource := rawContextSource{
			SourceName: f,
			RawContext: c,
		}
		sources = append(sources, rawSource)
	}
	rawContext := rawDirContext{
		RawSources: sources,
	}
	return rawContext.toCompilerInputContext()
}

func (e *Reader) decodeFile(file string) (rawFileContext, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return rawFileContext{}, err
	}
	c, err := e.decoder.decode(data)
	if err != nil {
		return rawFileContext{}, err
	}
	return c, nil
}
