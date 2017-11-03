package config

import (
	"errors"
	"fmt"

	"sort"

	. "github.com/dankraw/ssh-aliases/domain"
)

type HostsWithConfigs struct {
	Hosts         []Host        `hcl:"host"`
	RawSSHConfigs RawSSHConfigs `hcl:"config"`
}

type Host struct {
	Name        string      `hcl:",key"`
	Hostname    string      `hcl:"hostname"`
	Alias       string      `hcl:"alias"`
	ConfigOrRef ConfigOrRef `hcl:"config"`
}

type ConfigOrRef interface{}

type RawSSHConfigs map[string]interface{}

func (c *HostsWithConfigs) namedConfigs() (NamedConfigs, error) {
	namedConfigsMap := NamedConfigs{}
	for name, r := range c.RawSSHConfigs {
		if _, err := namedConfigsMap[name]; err {
			return NamedConfigs{}, errors.New(fmt.Sprintf("Duplicate config with name %v", name))
		}
		h := HostConfig{}
		if m, ok := r.([]map[string]interface{}); ok {
			for _, x := range m {
				for k, v := range x {
					if _, ok := h[k]; ok {
						return NamedConfigs{}, errors.New(fmt.Sprintf("Duplicate config entry `%v` in host `%v`", k, name))
					}
					h[k] = v
				}
			}
		} else {
			return NamedConfigs{}, errors.New(fmt.Sprintf("Config `%v` is not a key-value map", name))
		}
		namedConfigsMap[name] = h
	}
	return namedConfigsMap, nil
}

type NamedConfigs map[string]HostConfig

type HostConfig map[string]interface{}

func (c *NamedConfigs) toNamedHostConfigEntries() NamedHostConfigEntries {
	entries := NamedHostConfigEntries{}
	for n, config := range *c {
		entries[n] = config.toSortedHostConfigEntries()
	}
	return entries
}

type NamedHostConfigEntries map[string]HostConfigEntries

var sanitizer = NewSanitizer()

func (c *HostConfig) toSortedHostConfigEntries() HostConfigEntries {
	entries := []HostConfigEntry{}
	for k, v := range *c {
		entries = append(entries, HostConfigEntry{sanitizer.Sanitize(k), v})
	}
	sort.Sort(ByHostConfigEntryKey(entries))
	return entries
}

func (c *HostsWithConfigs) ToHostConfigInputs() ([]HostConfigInput, error) {
	namedConfigs, err := c.namedConfigs()
	if err != nil {
		return nil, err
	}
	namedConfigEntries := namedConfigs.toNamedHostConfigEntries()
	inputs := []HostConfigInput{}

	aliases := map[string]Host{}
	for _, a := range c.Hosts {
		if _, ok := aliases[a.Name]; ok {
			return nil, errors.New(fmt.Sprintf("Duplicate host `%v`", a.Name))
		}
		aliases[a.Name] = a

		input := HostConfigInput{
			AliasName:       a.Name,
			HostnamePattern: a.Hostname,
			AliasTemplate:   a.Alias,
		}
		if configName, ok := a.ConfigOrRef.(string); ok {
			if named, ok := namedConfigEntries[configName]; ok {
				input.HostConfig = named
			} else {
				return nil, errors.New(fmt.Sprintf("No config `%v` found (used by host `%v`)",
					configName, a.Name))
			}
		} else if m, ok := a.ConfigOrRef.([]map[string]interface{}); ok {
			h := HostConfig{}
			for _, x := range m {
				for k, v := range x {
					if _, ok := h[k]; ok {
						return nil, errors.New(fmt.Sprintf("Duplicate config property `%v` for host `%v`", k, a.Name))
					}
					h[k] = v
				}
			}
			input.HostConfig = h.toSortedHostConfigEntries()
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func (c *HostsWithConfigs) Merge(config HostsWithConfigs) error {
	c.Hosts = append(c.Hosts, config.Hosts...)
	if c.RawSSHConfigs == nil {
		c.RawSSHConfigs = RawSSHConfigs{}
	}
	for n, r := range config.RawSSHConfigs {
		if _, ok := c.RawSSHConfigs[n]; ok {
			return errors.New(fmt.Sprintf("Duplicate config `%s`", n))
		}
		c.RawSSHConfigs[n] = r
	}
	return nil
}
