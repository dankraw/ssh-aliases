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
	Hosts      []host       `hcl:"host"`
	RawConfigs rawConfigs   `hcl:"config"`
	Variables  rawVariables `hcl:"var"`
}

type rawVariables map[string]interface{}

type variablesMap map[string]string

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

func (vals *variablesMap) applyTo(str string) string {
	if strings.Contains(str, "${") {
		for k, v := range *vals {
			str = strings.Replace(str, fmt.Sprintf("${%s}", k), v, -1)
		}
	}
	return str
}

func interpolatedConfigProps(variables *variablesMap, input []map[string]interface{}) configProps {
	h := configProps{}
	for _, x := range input {
		for k, v := range x {
			if vStr, ok := v.(string); ok {
				h[k] = variables.applyTo(vStr)
			} else {
				h[k] = v
			}
		}
	}
	return h
}

func (c *rawDirContext) getConfigPropertiesMap(variables *variablesMap) configPropertiesMap {
	propsMap := configPropertiesMap{}
	for _, s := range c.RawSources {
		for name, r := range s.RawContext.RawConfigs {
			interpolated := interpolatedConfigProps(variables, r)
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

func (c *rawFileContext) toExpandingHostConfigs(variables *variablesMap, propsMap *configPropertiesMap) ([]compiler.ExpandingHostConfig, error) {
	configsMap := *propsMap
	inputs := []compiler.ExpandingHostConfig{}

	for _, a := range c.Hosts {
		input := compiler.ExpandingHostConfig{
			AliasName:       a.Name,
			HostnamePattern: variables.applyTo(a.Hostname),
			AliasTemplate:   variables.applyTo(a.Alias),
		}
		if configName, ok := a.RawConfigOrRef.(string); ok {
			if named, ok := configsMap[configName]; ok {
				input.Config = named
			} else {
				return nil, fmt.Errorf("no config `%v` found (used by host `%v`)",
					configName, a.Name)
			}
		} else if m, ok := a.RawConfigOrRef.([]map[string]interface{}); ok {
			interpolated := interpolatedConfigProps(variables, m)
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
	variables, err := c.getNormalizedVariables()
	if err != nil {
		return compiler.InputContext{}, err
	}
	propsMap := c.getConfigPropertiesMap(&variables)
	if err != nil {
		return compiler.InputContext{}, err
	}
	var sources []compiler.ContextSource
	for _, s := range c.RawSources {
		expandingHostConfigs, err := s.RawContext.toExpandingHostConfigs(&variables, &propsMap)
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

func (c *rawDirContext) getNormalizedVariables() (variablesMap, error) {
	variables := variablesMap{}
	for _, s := range c.RawSources {
		for k, v := range s.RawContext.Variables {
			for key, variable := range c.expandVariable(k, v) {
				if _, contains := variables[key]; contains {
					return nil, fmt.Errorf("variable redeclaration: %v", key)
				}
				variables[key] = variable
			}
		}
	}
	return variables, nil
}

func (c *rawDirContext) expandVariable(key string, variable interface{}) map[string]string {
	expanded := map[string]string{}
	if arr, ok := variable.([]map[string]interface{}); ok {
		for _, m := range arr {
			for k, v := range m {
				for ek, ev := range c.expandVariable(k, v) {
					expanded[fmt.Sprintf("%s.%s", key, ek)] = ev
				}
			}
		}
	} else {
		expanded[key] = fmt.Sprintf("%v", variable)
	}
	return expanded
}
