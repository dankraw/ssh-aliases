package config

import (
	"io/ioutil"

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
func (e *Reader) ReadConfigs(dir string) ([]compiler.ExpandingHostConfig, error) {
	files, err := e.scanner.ScanDirectory(dir)
	if err != nil {
		return nil, err
	}
	var parsed []rawConfigContext
	for _, f := range files {
		c, err := e.decodeFile(f)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, c)
	}
	merged, err := mergeRawConfigContexts(parsed...)
	if err != nil {
		return nil, err

	}
	configs, err := merged.toExpandingHostConfigs()
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// ReadConfig processes the input file and returns inputs for ssh-aliases compiler
func (e *Reader) ReadConfig(file string) ([]compiler.ExpandingHostConfig, error) {
	config, err := e.decodeFile(file)
	if err != nil {
		return nil, err
	}
	return config.toExpandingHostConfigs()
}

func (e *Reader) decodeFile(file string) (rawConfigContext, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return rawConfigContext{}, err
	}
	c, err := e.decoder.decode(data)
	if err != nil {
		return rawConfigContext{}, err
	}
	return c, nil
}
