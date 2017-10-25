package config

import (
	"github.com/hashicorp/hcl"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) decode(input []byte) (Config, error) {
	config := Config{}
	file, err := hcl.ParseBytes(input)
	if err != nil {
		return Config{}, err
	}
	err = hcl.DecodeObject(&config, file)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
