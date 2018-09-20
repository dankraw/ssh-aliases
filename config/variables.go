package config

import "fmt"

func normalizedVariables(sources []rawContextSource) (variablesMap, error) {
	variables := variablesMap{}
	for _, s := range sources {
		for k, v := range s.RawContext.Variables {
			for key, variable := range expandVariable(k, v) {
				if _, contains := variables[key]; contains {
					return nil, fmt.Errorf("error in `%s`: variable redeclaration: `%v`", s.SourceName, key)
				}
				variables[key] = variable
			}
		}
	}
	return variables, nil
}

func expandVariable(key string, variable interface{}) map[string]string {
	expanded := map[string]string{}
	if arr, ok := variable.([]map[string]interface{}); ok {
		for _, m := range arr {
			for k, v := range m {
				for ek, ev := range expandVariable(k, v) {
					expanded[fmt.Sprintf("%s.%s", key, ek)] = ev
				}
			}
		}
	} else {
		expanded[key] = fmt.Sprintf("%v", variable)
	}
	return expanded
}
