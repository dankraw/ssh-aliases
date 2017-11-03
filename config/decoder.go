package config

import (
	"github.com/hashicorp/hcl"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(input []byte) (HostsWithConfigs, error) {
	config := HostsWithConfigs{}
	file, err := hcl.ParseBytes(input)
	if err != nil {
		return HostsWithConfigs{}, err
	}
	err = hcl.DecodeObject(&config, file)
	if err != nil {
		return HostsWithConfigs{}, err
	}
	return config, nil
}
