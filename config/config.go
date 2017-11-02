package config

import (
	"errors"
	"fmt"

	. "github.com/dankraw/ssh-aliases/domain"
)

type Config struct {
	Aliases       []Alias       `hcl:"alias"`
	RawSSHConfigs RawSSHConfigs `hcl:"ssh_config"`
}

type Alias struct {
	Name          string     `hcl:",key"`
	Pattern       string     `hcl:"pattern"`
	Template      string     `hcl:"template"`
	SSHConfigName string     `hcl:"ssh_config_name"`
	SSHConfig     HostConfig `hcl:"ssh_config"`
}

type RawSSHConfigs map[string]interface{}

type NamedConfigs map[string]HostConfig

func (c *Config) NamedConfigs() (NamedConfigs, error) {
	namedConfigsMap := NamedConfigs{}
	for name, r := range c.RawSSHConfigs {
		if _, err := namedConfigsMap[name]; err {
			return NamedConfigs{}, errors.New(fmt.Sprintf("Duplicate ssh-config with name %v", name))
		}
		namedConfigsMap[name] = r.([]map[string]interface{})[0]
	}
	return namedConfigsMap, nil
}

func (c *Config) ToHostConfigInputs() ([]HostConfigInput, error) {
	namedConfigs, err := c.NamedConfigs()
	if err != nil {
		return nil, err
	}
	inputs := []HostConfigInput{}

	aliases := map[string]Alias{}
	for _, a := range c.Aliases {
		if _, ok := aliases[a.Name]; ok {
			return nil, errors.New(fmt.Sprintf("Duplicate alias with name %v", a.Name))
		}
		aliases[a.Name] = a
		input := HostConfigInput{
			AliasName:       a.Name,
			HostnamePattern: a.Pattern,
			AliasTemplate:   a.Template,
		}
		if a.SSHConfigName != "" {
			if named, ok := namedConfigs[a.SSHConfigName]; ok {
				input.HostConfig = named
			} else {
				return nil, errors.New(fmt.Sprintf("No ssh-config named %v found (used by %v alias)",
					a.SSHConfigName, a.Name))
			}
		} else {
			input.HostConfig = a.SSHConfig
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func (c *Config) Merge(config Config) error {
	c.Aliases = append(c.Aliases, config.Aliases...)
	if c.RawSSHConfigs == nil {
		c.RawSSHConfigs = RawSSHConfigs{}
	}
	for k, r := range config.RawSSHConfigs {
		if _, ok := c.RawSSHConfigs[k]; ok {
			return errors.New(fmt.Sprintf("Duplicate ssh-config with name %s", k))
		}
		c.RawSSHConfigs[k] = r
	}
	return nil
}
