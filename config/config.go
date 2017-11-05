package config

import (
	"errors"
	"fmt"

	"sort"

	. "github.com/dankraw/ssh-aliases/domain"
)

type HostsWithConfigs struct {
	Hosts      []Host     `hcl:"host"`
	RawConfigs RawConfigs `hcl:"config"`
}

type Host struct {
	Name           string         `hcl:",key"`
	Hostname       string         `hcl:"hostname"`
	Alias          string         `hcl:"alias"`
	RawConfigOrRef RawConfigOrRef `hcl:"config"`
}

type RawConfigOrRef interface{}

type RawConfigs map[string]interface{}

func (c *HostsWithConfigs) hostConfigsMap() (HostConfigsMap, error) {
	namedConfigsMap := HostConfigsMap{}
	for name, r := range c.RawConfigs {
		if _, err := namedConfigsMap[name]; err {
			return HostConfigsMap{}, errors.New(fmt.Sprintf("Duplicate config with name %v", name))
		}
		h := HostConfig{}
		if m, ok := r.([]map[string]interface{}); ok {
			for _, x := range m {
				for k, v := range x {
					if _, ok := h[k]; ok {
						return HostConfigsMap{}, errors.New(fmt.Sprintf("Duplicate config entry `%v` in host `%v`", k, name))
					}
					h[k] = v
				}
			}
		} else {
			return HostConfigsMap{}, errors.New(fmt.Sprintf("Config `%v` is not a key-value map", name))
		}
		namedConfigsMap[name] = h
	}
	return namedConfigsMap, nil
}

type HostConfigsMap map[string]HostConfig

type HostConfig map[string]interface{}

func (c *HostConfigsMap) toHostConfigEntriesMap() HostConfigEntriesMap {
	entries := HostConfigEntriesMap{}
	for n, config := range *c {
		entries[n] = config.toSortedHostConfigEntries()
	}
	return entries
}

type HostConfigEntriesMap map[string]HostConfigEntries

var sanitizer = NewSanitizer()

func (c *HostConfig) toSortedHostConfigEntries() HostConfigEntries {
	var entries []HostConfigEntry
	for k, v := range *c {
		entries = append(entries, HostConfigEntry{sanitizer.Sanitize(k), v})
	}
	sort.Sort(ByHostConfigEntryKey(entries))
	return entries
}

func (c *HostsWithConfigs) ToHostConfigInputs() ([]HostConfigInput, error) {
	configsMap, err := c.hostConfigsMap()
	if err != nil {
		return nil, err
	}
	namedConfigEntries := configsMap.toHostConfigEntriesMap()
	var inputs []HostConfigInput

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
		if configName, ok := a.RawConfigOrRef.(string); ok {
			if named, ok := namedConfigEntries[configName]; ok {
				input.HostConfig = named
			} else {
				return nil, errors.New(fmt.Sprintf("No config `%v` found (used by host `%v`)",
					configName, a.Name))
			}
		} else if m, ok := a.RawConfigOrRef.([]map[string]interface{}); ok {
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
	if c.RawConfigs == nil {
		c.RawConfigs = RawConfigs{}
	}
	for n, r := range config.RawConfigs {
		if _, ok := c.RawConfigs[n]; ok {
			return errors.New(fmt.Sprintf("Duplicate config `%s`", n))
		}
		c.RawConfigs[n] = r
	}
	return nil
}
