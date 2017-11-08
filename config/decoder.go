package config

import (
	"github.com/hashicorp/hcl"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(input []byte) (RawConfigContext, error) {
	config := RawConfigContext{}
	file, err := hcl.ParseBytes(input)
	if err != nil {
		return RawConfigContext{}, err
	}
	err = hcl.DecodeObject(&config, file)
	if err != nil {
		return RawConfigContext{}, err
	}
	return config, nil
}
