package config

import (
	"fmt"
	"regexp"
	"strings"
)

var variableRegexp = regexp.MustCompile("\\${([a-zA-Z0-9-\\.]+)}")

type configProps map[string]interface{}

type variablesMap map[string]string

func interpolatedConfigProps(variables variablesMap, rawConfig []map[string]interface{}) (configProps, error) {
	h := configProps{}
	for _, x := range rawConfig {
		for k, v := range x {
			if vStr, ok := v.(string); ok {
				interpolated, err := applyVariablesToString(vStr, variables)
				if err != nil {
					return nil, fmt.Errorf("could not compile config property `%s`: %s", k, err.Error())
				}
				h[k] = interpolated
			} else {
				h[k] = v
			}
		}
	}
	return h, nil
}

func applyVariablesToString(str string, vals variablesMap) (string, error) {
	if strings.Contains(str, "${") {
		for k, v := range vals {
			str = strings.Replace(str, fmt.Sprintf("${%s}", k), v, -1)
		}
		matches := variableRegexp.FindAllStringSubmatch(str, -1)
		if matches != nil {
			return "", fmt.Errorf("variable `%s` not defined", matches[0][1])
		}
	}
	return str, nil
}

const importConfigKey = "_import"

func (c configProps) evaluateConfigImports(propsMap map[string]configProps, evaluatedImports []string) (configProps, error) {
	evaluated := configProps{}
	for key, value := range c {
		if key == importConfigKey {
			if importedStr, ok := value.(string); ok {
				if contains(evaluatedImports, importedStr) {
					return nil, fmt.Errorf("circular import in configs (config imports chain: `%s` -> `%s`)", strings.Join(evaluatedImports, " -> "), importedStr)
				}
				evaluatedImports = append(evaluatedImports, importedStr)
				if imported, ok := propsMap[importedStr]; ok {
					evaluatedImport, err := imported.evaluateConfigImports(propsMap, evaluatedImports)
					if err != nil {
						return nil, err
					}
					for k, v := range evaluatedImport {
						evaluated[k] = v
					}
				} else {
					return nil, fmt.Errorf("trying to import `%s`, but such config does not exist", importedStr)
				}
			} else {
				return nil, fmt.Errorf("config import statement has invalid value: `%v`", value)
			}
		} else {
			evaluated[key] = value
		}
	}
	return evaluated, nil
}

func contains(slice []string, element string) bool {
	for _, e := range slice {
		if element == e {
			return true
		}
	}
	return false
}
