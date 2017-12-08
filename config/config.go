package config

import (
	"fmt"
	"strings"

	"sort"

	"github.com/dankraw/ssh-aliases/compiler"
)

type rawDirContext struct {
	RawSources []rawContextSource
}

type rawContextSource struct {
	SourceName string
	RawContext rawFileContext
}

type rawFileContext struct {
	Hosts      []host     `hcl:"host"`
	RawConfigs rawConfigs `hcl:"config"`
	Values     rawValues  `hcl:"values"`
}

type rawValues map[string]interface{}

type valuesMap map[string]string

type host struct {
	Name           string         `hcl:",key"`
	Hostname       string         `hcl:"hostname"`
	Alias          string         `hcl:"alias"`
	RawConfigOrRef rawConfigOrRef `hcl:"config"`
}

type rawConfigOrRef interface{}

type rawConfigs map[string]rawConfig

type rawConfig []map[string]interface{}

type configProps map[string]interface{}

func (vals *valuesMap) applyTo(str string) string {
	if strings.Contains(str, "${") {
		for k, v := range *vals {
			str = strings.Replace(str, fmt.Sprintf("${%s}", k), v, -1)
		}
	}
	return str
}

func interpolatedConfigProps(values *valuesMap, input []map[string]interface{}) configProps {
	h := configProps{}
	for _, x := range input {
		for k, v := range x {
			if vStr, ok := v.(string); ok {
				h[k] = values.applyTo(vStr)
			} else {
				h[k] = v
			}
		}
	}
	return h
}

func (c *rawDirContext) getConfigPropertiesMap(values *valuesMap) configPropertiesMap {
	propsMap := configPropertiesMap{}
	for _, s := range c.RawSources {
		for name, r := range s.RawContext.RawConfigs {
			interpolated := interpolatedConfigProps(values, r)
			propsMap[name] = interpolated.toSortedCompilerProperties()
		}
	}
	return propsMap
}

type configPropertiesMap map[string]compiler.ConfigProperties

var sanitizer = newKeywordSanitizer()

func (c *configProps) toSortedCompilerProperties() compiler.ConfigProperties {
	var entries []compiler.ConfigProperty
	for k, v := range *c {
		entries = append(entries, compiler.ConfigProperty{Key: sanitizer.sanitize(k), Value: v})
	}
	sort.Sort(compiler.ByConfigPropertyKey(entries))
	return entries
}

func (c *rawFileContext) toExpandingHostConfigs(values *valuesMap, propsMap *configPropertiesMap) ([]compiler.ExpandingHostConfig, error) {
	configsMap := *propsMap
	inputs := []compiler.ExpandingHostConfig{}

	for _, a := range c.Hosts {
		input := compiler.ExpandingHostConfig{
			AliasName:       a.Name,
			HostnamePattern: values.applyTo(a.Hostname),
			AliasTemplate:   values.applyTo(a.Alias),
		}
		if configName, ok := a.RawConfigOrRef.(string); ok {
			if named, ok := configsMap[configName]; ok {
				input.Config = named
			} else {
				return nil, fmt.Errorf("no config `%v` found (used by host `%v`)",
					configName, a.Name)
			}
		} else if m, ok := a.RawConfigOrRef.([]map[string]interface{}); ok {
			// @TODO duplication on parsing configs?
			interpolated := interpolatedConfigProps(values, m)
			input.Config = interpolated.toSortedCompilerProperties()
		} else {
			return nil, fmt.Errorf("invalid config definition for host `%v`", a.Name)
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func (c *rawDirContext) toCompilerInputContext() (compiler.InputContext, error) {
	err := c.validateHosts()
	if err != nil {
		return compiler.InputContext{}, err
	}
	values, err := c.getNormalizedValues()
	if err != nil {
		return compiler.InputContext{}, err
	}
	propsMap := c.getConfigPropertiesMap(&values)
	if err != nil {
		return compiler.InputContext{}, err
	}
	var sources []compiler.ContextSource
	for _, s := range c.RawSources {
		expandingHostConfigs, err := s.RawContext.toExpandingHostConfigs(&values, &propsMap)
		if err != nil {
			return compiler.InputContext{}, err
		}
		sources = append(sources, compiler.ContextSource{
			SourceName: s.SourceName,
			Hosts:      expandingHostConfigs,
		})
	}
	return compiler.InputContext{
		Sources: sources,
	}, nil
}

func (c *rawDirContext) validateHosts() error {
	hosts := make(map[string]struct{})
	var exists struct{}
	for _, s := range c.RawSources {
		for _, h := range s.RawContext.Hosts {
			if _, contains := hosts[h.Name]; contains {
				return fmt.Errorf("duplicate host `%v`", h.Name)
			}
			hosts[h.Name] = exists
		}
	}
	return nil
}

func (c *rawDirContext) getNormalizedValues() (valuesMap, error) {
	vals := valuesMap{}
	for _, s := range c.RawSources {
		for k, v := range s.RawContext.Values {
			for key, value := range c.expandValue(k, v) {
				if _, contains := vals[key]; contains {
					return nil, fmt.Errorf("value redeclaration: %v", key)
				}
				vals[key] = value
			}
		}
	}
	return vals, nil
}

func (c *rawDirContext) expandValue(key string, value interface{}) map[string]string {
	expanded := map[string]string{}
	if arr, ok := value.([]map[string]interface{}); ok {
		for _, m := range arr {
			for k, v := range m {
				for ek, ev := range c.expandValue(k, v) {
					expanded[fmt.Sprintf("%s.%s", key, ek)] = ev
				}
			}
		}
	} else {
		expanded[key] = fmt.Sprintf("%v", value)
	}
	return expanded
}
