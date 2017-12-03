package config

import (
	"github.com/hashicorp/hcl"
)

type decoder struct{}

func newDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) decode(input []byte) (rawFileContext, error) {
	config := rawFileContext{}
	file, err := hcl.ParseBytes(input)
	if err != nil {
		return rawFileContext{}, err
	}
	err = hcl.DecodeObject(&config, file)
	if err != nil {
		return rawFileContext{}, err
	}
	return config, nil
}
