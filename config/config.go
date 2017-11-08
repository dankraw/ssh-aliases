package config

import (
	"errors"
	"fmt"

	"sort"

	"github.com/dankraw/ssh-aliases/compiler"
)

type RawConfigContext struct {
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

type RawConfigs map[string]RawConfig

type RawConfig []Config

type Config map[string]interface{}

func (c *RawConfigContext) toConfigPropertiesMap() (configPropertiesMap, error) {
	propsMap := configPropertiesMap{}
	for name, r := range c.RawConfigs {
		if _, exists := propsMap[name]; exists {
			return configPropertiesMap{}, errors.New(fmt.Sprintf("Duplicate config with name %v", name))
		}
		h := Config{}
		for _, x := range r {
			for k, v := range x {
				if _, ok := h[k]; ok {
					return configPropertiesMap{}, errors.New(fmt.Sprintf("Duplicate config entry `%v` in host `%v`", k, name))
				}
				h[k] = v
			}
		}
		propsMap[name] = h.toSortedProperties()
	}
	return propsMap, nil
}

type configPropertiesMap map[string]compiler.ConfigProperties

var sanitizer = NewKeywordSanitizer()

func (c *Config) toSortedProperties() compiler.ConfigProperties {
	var entries []compiler.ConfigProperty
	for k, v := range *c {
		entries = append(entries, compiler.ConfigProperty{sanitizer.Sanitize(k), v})
	}
	sort.Sort(compiler.ByConfigPropertyKey(entries))
	return entries
}

func (c *RawConfigContext) ToExpandingHostConfigs() ([]compiler.ExpandingHostConfig, error) {
	configsMap, err := c.toConfigPropertiesMap()
	if err != nil {
		return nil, err
	}
	var inputs []compiler.ExpandingHostConfig

	aliases := map[string]Host{}
	for _, a := range c.Hosts {
		if _, ok := aliases[a.Name]; ok {
			return nil, errors.New(fmt.Sprintf("Duplicate host `%v`", a.Name))
		}
		aliases[a.Name] = a

		input := compiler.ExpandingHostConfig{
			AliasName:       a.Name,
			HostnamePattern: a.Hostname,
			AliasTemplate:   a.Alias,
		}
		if configName, ok := a.RawConfigOrRef.(string); ok {
			if named, ok := configsMap[configName]; ok {
				input.Config = named
			} else {
				return nil, errors.New(fmt.Sprintf("No config `%v` found (used by host `%v`)",
					configName, a.Name))
			}
		} else if m, ok := a.RawConfigOrRef.([]map[string]interface{}); ok {
			h := Config{}
			for _, x := range m {
				for k, v := range x {
					if _, ok := h[k]; ok {
						return nil, errors.New(fmt.Sprintf("Duplicate config property `%v` for host `%v`", k, a.Name))
					}
					h[k] = v
				}
			}
			input.Config = h.toSortedProperties()
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid config definition for host `%v`", a.Name))
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func (c *RawConfigContext) Merge(config RawConfigContext) error {
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
