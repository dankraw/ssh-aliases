package config

import (
	"github.com/hashicorp/hcl"
)

type Config struct {
	Aliases    []Alias     `hcl:"alias"`
	SSHConfigs []SSHConfig `hcl:"ssh_config"`
}

type Alias struct {
	Name          string   `hcl:",key"`
	Patterns      []string `hcl:"patterns"`
	RegExp        string   `hcl:"regexp"`
	Template      string   `hcl:"template"`
	SSHConfigName string   `hcl:"ssh_config_name"`
}

type SSHConfig struct {
	Name         string `hcl:",key"`
	IdentityFile string `hcl:"identity_file"`
	Port int `hcl:"port"`
}

type Decoder struct {}

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