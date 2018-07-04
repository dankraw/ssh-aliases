package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dankraw/ssh-aliases/compiler"
)

type rawContextSource struct {
	SourceName string
	RawContext rawFileContext
}

func compilerInputContext(sources []rawContextSource) (compiler.InputContext, error) {
	err := validateHosts(sources)
	if err != nil {
		return compiler.InputContext{}, err
	}
	variables, err := normalizedVariables(sources)
	if err != nil {
		return compiler.InputContext{}, err
	}
	namedProps, err := getNamedConfigProps(sources, variables)
	if err != nil {
		return compiler.InputContext{}, err
	}
	var ctxSources []compiler.ContextSource
	for _, s := range sources {
		expandingHostConfigs, err := expandingHostConfigs(s.RawContext, variables, namedProps)
		if err != nil {
			return compiler.InputContext{}, err
		}
		ctxSources = append(ctxSources, compiler.ContextSource{
			SourceName: s.SourceName,
			Hosts:      expandingHostConfigs,
		})
	}
	return compiler.InputContext{
		Sources: ctxSources,
	}, nil
}

func validateHosts(sources []rawContextSource) error {
	hosts := make(map[string]struct{})
	var exists struct{}
	for _, s := range sources {
		for _, h := range s.RawContext.Hosts {
			if strings.TrimSpace(h.Alias) == "" {
				return fmt.Errorf("host definition `%v` contains no valid alias property", h.Name)
			}
			if _, contains := hosts[h.Name]; contains {
				return fmt.Errorf("duplicate host `%v`", h.Name)
			}
			hosts[h.Name] = exists
		}
	}
	return nil
}

func getNamedConfigProps(sources []rawContextSource, variables variablesMap) (map[string]configProps, error) {
	propsMap := map[string]configProps{}
	for _, s := range sources {
		for name, r := range s.RawContext.RawConfigs {
			propsMap[name] = interpolatedConfigProps(variables, r)
		}
	}
	evaluated := map[string]configProps{}
	for name, props := range propsMap {
		evaluatedConfig, err := props.evaluateConfigImports(propsMap, make([]string, 0))
		if err != nil {
			return nil, err
		}
		evaluated[name] = evaluatedConfig
	}
	return evaluated, nil
}

func expandingHostConfigs(fileCtx rawFileContext, variables variablesMap, propsMap map[string]configProps) ([]compiler.ExpandingHostConfig, error) {
	configsMap := propsMap
	inputs := []compiler.ExpandingHostConfig{}

	for _, a := range fileCtx.Hosts {
		config := compiler.ConfigProperties{}
		if configName, ok := a.RawConfigOrRef.(string); ok {
			if named, ok := configsMap[configName]; ok {
				config = sortedCompilerProperties(named)
			} else {
				return nil, fmt.Errorf("no config `%v` found (used by host `%v`)",
					configName, a.Name)
			}
		} else if m, ok := a.RawConfigOrRef.([]map[string]interface{}); ok {
			interpolated := interpolatedConfigProps(variables, m)
			evaluated, err := interpolated.evaluateConfigImports(configsMap, make([]string, 0))
			if err != nil {
				return nil, err
			}
			config = sortedCompilerProperties(evaluated)
		} else if a.RawConfigOrRef == nil {
			if strings.TrimSpace(a.Hostname) == "" {
				return nil, fmt.Errorf("no config nor hostname specified for for host `%v`", a.Name)
			}
		} else {
			return nil, fmt.Errorf("invalid config definition for host `%v`", a.Name)
		}
		inputs = append(inputs, compiler.ExpandingHostConfig{
			AliasName:       a.Name,
			HostnamePattern: applyVariablesToString(a.Hostname, variables),
			AliasTemplate:   applyVariablesToString(a.Alias, variables),
			Config:          config,
		})
	}
	return inputs, nil
}

func sortedCompilerProperties(props configProps) compiler.ConfigProperties {
	var entries []compiler.ConfigProperty
	for k, v := range props {
		entries = append(entries, compiler.ConfigProperty{Key: sanitize(k), Value: v})
	}
	sort.Sort(compiler.ByConfigPropertyKey(entries))
	return entries
}
