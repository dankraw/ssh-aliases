package config

import (
	"github.com/hashicorp/hcl"
)

type decoder struct{}

func newDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) decode(input []byte) (rawConfigContext, error) {
	config := rawConfigContext{}
	file, err := hcl.ParseBytes(input)
	if err != nil {
		return rawConfigContext{}, err
	}
	err = hcl.DecodeObject(&config, file)
	if err != nil {
		return rawConfigContext{}, err
	}
	return config, nil
}
