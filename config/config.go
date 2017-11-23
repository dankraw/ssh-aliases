package config

import (
	"fmt"

	"sort"

	"github.com/dankraw/ssh-aliases/compiler"
)

type rawConfigContext struct {
	Hosts      []host     `hcl:"host"`
	RawConfigs rawConfigs `hcl:"config"`
}

type host struct {
	Name           string         `hcl:",key"`
	Hostname       string         `hcl:"hostname"`
	Alias          string         `hcl:"alias"`
	RawConfigOrRef rawConfigOrRef `hcl:"config"`
}

type rawConfigOrRef interface{}

type rawConfigs map[string]rawConfig

type rawConfig []configProps

type configProps map[string]interface{}

func (c *rawConfigContext) toConfigPropertiesMap() (configPropertiesMap, error) {
	propsMap := configPropertiesMap{}
	for name, r := range c.RawConfigs {
		if _, exists := propsMap[name]; exists {
			return configPropertiesMap{}, fmt.Errorf("duplicate config with name %v", name)
		}
		h := configProps{}
		for _, x := range r {
			for k, v := range x {
				if _, ok := h[k]; ok {
					return configPropertiesMap{}, fmt.Errorf("duplicate config entry `%v` in host `%v`", k, name)
				}
				h[k] = v
			}
		}
		propsMap[name] = h.toSortedProperties()
	}
	return propsMap, nil
}

type configPropertiesMap map[string]compiler.ConfigProperties

var sanitizer = newKeywordSanitizer()

func (c *configProps) toSortedProperties() compiler.ConfigProperties {
	var entries []compiler.ConfigProperty
	for k, v := range *c {
		entries = append(entries, compiler.ConfigProperty{Key: sanitizer.sanitize(k), Value: v})
	}
	sort.Sort(compiler.ByConfigPropertyKey(entries))
	return entries
}

func (c *rawConfigContext) toExpandingHostConfigs() ([]compiler.ExpandingHostConfig, error) {
	configsMap, err := c.toConfigPropertiesMap()
	if err != nil {
		return nil, err
	}
	var inputs []compiler.ExpandingHostConfig

	aliases := map[string]host{}
	for _, a := range c.Hosts {
		if _, ok := aliases[a.Name]; ok {
			return nil, fmt.Errorf("duplicate host `%v`", a.Name)
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
				return nil, fmt.Errorf("no config `%v` found (used by host `%v`)",
					configName, a.Name)
			}
		} else if m, ok := a.RawConfigOrRef.([]map[string]interface{}); ok {
			h := configProps{}
			for _, x := range m {
				for k, v := range x {
					if _, ok := h[k]; ok {
						return nil, fmt.Errorf("duplicate config property `%v` for host `%v`", k, a.Name)
					}
					h[k] = v
				}
			}
			input.Config = h.toSortedProperties()
		} else {
			return nil, fmt.Errorf("invalid config definition for host `%v`", a.Name)
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func mergeRawConfigContexts(contexts ...rawConfigContext) (rawConfigContext, error) {
	m := rawConfigContext{
		Hosts:      []host{},
		RawConfigs: rawConfigs{},
	}
	for _, c := range contexts {
		m.Hosts = append(m.Hosts, c.Hosts...)
		for n, r := range c.RawConfigs {
			if _, ok := m.RawConfigs[n]; ok {
				return rawConfigContext{}, fmt.Errorf("duplicate config `%s`", n)
			}
			m.RawConfigs[n] = r
		}
	}
	return m, nil
}
