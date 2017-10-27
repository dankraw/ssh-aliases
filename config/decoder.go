package config

import (
	"github.com/hashicorp/hcl"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(input []byte) (Config, error) {
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
