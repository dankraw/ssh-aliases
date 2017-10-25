package config

import (
	"errors"
	"fmt"
	. "github.com/dankraw/ssh-aliases/domain"
)

type Config struct {
	Aliases    []Alias     `hcl:"alias"`
	SSHConfigs []SSHConfig `hcl:"ssh_config"`
}

type Alias struct {
	Name          string            `hcl:",key"`
	Pattern       string            `hcl:"pattern"`
	Template      string            `hcl:"template"`
	SSHConfigName string            `hcl:"ssh_config_name"`
	SSHConfig     EmbeddedSSHConfig `hcl:"ssh_config"`
}

type EmbeddedSSHConfig struct {
	IdentityFile string `hcl:"identity_file"`
	Port         int    `hcl:"port"`
}

func (c *EmbeddedSSHConfig) ToHostConfig() HostConfig {
	return HostConfig{
		IdentityFile: c.IdentityFile,
		Port:         c.Port,
	}
}

type SSHConfig struct {
	Name         string `hcl:",key"`
	IdentityFile string `hcl:"identity_file"`
	Port         int    `hcl:"port"`
}

func (c *SSHConfig) ToHostConfig() HostConfig {
	return HostConfig{
		IdentityFile: c.IdentityFile,
		Port:         c.Port,
	}
}

func (c *Config) ToHostConfigInputs() ([]HostConfigInput, error) {
	inputs := []HostConfigInput{}

	namedConfigsMap := map[string]HostConfig{}
	for _, named := range c.SSHConfigs {
		namedConfigsMap[named.Name] = named.ToHostConfig()
	}
	for _, a := range c.Aliases {
		input := HostConfigInput{
			HostnamePattern: a.Pattern,
			AliasTemplate:   a.Template,
		}
		if a.SSHConfigName != "" {
			if named, ok := namedConfigsMap[a.SSHConfigName]; ok {
				input.HostConfig = named
			} else {
				return nil, errors.New(fmt.Sprintf("No ssh-config named %v found (used by %v alias)",
					a.SSHConfigName, a.Name))
			}
		} else {
			input.HostConfig = a.SSHConfig.ToHostConfig()
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}
