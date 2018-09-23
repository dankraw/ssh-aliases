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
			return compiler.InputContext{}, fmt.Errorf("error in `%s`: %s", s.SourceName, err.Error())
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
				return fmt.Errorf("error in `%s`: invalid `%s` host definition: empty alias", s.SourceName, h.Name)
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
	configToSourceMap := map[string]string{}
	propsMap := map[string]configProps{}
	for _, s := range sources {
		for name, r := range s.RawContext.RawConfigs {
			configToSourceMap[name] = s.SourceName
			interpolated, err := interpolatedConfigProps(variables, r)
			if err != nil {
				return nil, fmt.Errorf("error in `%s`: invalid `%s` config definition: %s", s.SourceName, name, err.Error())
			}
			propsMap[name] = interpolated
		}
	}
	evaluated := map[string]configProps{}
	for name, props := range propsMap {
		evaluatedImports := make([]string, 0)
		evaluatedConfig, err := props.evaluateConfigImports(propsMap, &evaluatedImports)
		if err != nil {
			return nil, fmt.Errorf("error in `%s`: invalid `%s` config definition: %s", configToSourceMap[name], name, err.Error())
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

		switch v := a.RawConfigOrRef.(type) {
		case string:
			if named, ok := configsMap[v]; ok {
				config = sortedCompilerProperties(named)
			} else {
				return nil, fmt.Errorf("error in `%s` host definition: no config `%s` found",
					a.Name, v)
			}
		case []map[string]interface{}:
			interpolated, err := interpolatedConfigProps(variables, v)
			if err != nil {
				return nil, fmt.Errorf("error in `%s` host definition: %s", a.Name, err.Error())
			}
			evaluatedImports := make([]string, 0)
			evaluated, err := interpolated.evaluateConfigImports(configsMap, &evaluatedImports)
			if err != nil {
				return nil, fmt.Errorf("error in `%s` host definition: %s", a.Name, err.Error())
			}
			config = sortedCompilerProperties(evaluated)
		case nil:
			if strings.TrimSpace(a.Hostname) == "" {
				return nil, fmt.Errorf("no config nor hostname specified for host `%v`", a.Name)
			}
		default:
			return nil, fmt.Errorf("invalid config definition for host `%v`", a.Name)
		}

		interpolatedHostname, err := applyVariablesToString(a.Hostname, variables)
		if err != nil {
			return nil, fmt.Errorf("error in hostname of `%s` host definition: %s", a.Name, err.Error())
		}
		interpolatedAlias, err := applyVariablesToString(a.Alias, variables)
		if err != nil {
			return nil, fmt.Errorf("error in alias of `%s` host definition: %s", a.Name, err.Error())
		}
		inputs = append(inputs, compiler.ExpandingHostConfig{
			AliasName:       a.Name,
			HostnamePattern: interpolatedHostname,
			AliasTemplate:   interpolatedAlias,
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
