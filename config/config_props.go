package config

import (
	"fmt"
	"regexp"
	"strings"
)

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

var variableRegexp = regexp.MustCompile(`\${([^}]+)}`)

func applyVariablesToString(str string, vals variablesMap) (string, error) {
	match := variableRegexp.FindStringSubmatchIndex(str)
	for match != nil {
		beginMatch := match[0]
		endMatch := match[1]
		beginIdx := match[2]
		endIdx := match[3]
		varName := str[beginIdx:endIdx]
		if value, ok := vals[varName]; ok {
			str = str[0:beginMatch] + value + str[endMatch:]
			match = variableRegexp.FindStringSubmatchIndex(str)
		} else {
			return "", fmt.Errorf("variable `%s` not defined", varName)
		}
	}
	return str, nil
}

const extendConfigKey = "_extend"

func (c configProps) evaluateConfigImports(propsMap map[string]configProps, evaluatedImports *[]string) (configProps, error) {
	if value, ok := c[extendConfigKey]; ok {
		evaluated := configProps{}
		if importedStr, ok := value.(string); ok {
			imported, err := importProps(importedStr, propsMap, evaluatedImports)
			if err != nil {
				return nil, err
			}
			for k, v := range imported {
				evaluated[k] = v
			}
		} else if importedArr, ok := value.([]interface{}); ok {
			for _, importedInterface := range importedArr {
				if importedStr, ok := importedInterface.(string); ok {

					// each import branch needs a copy of evaluated imports list
					evaluatedImportsBranch := make([]string, len(*evaluatedImports))
					copy(evaluatedImportsBranch, *evaluatedImports)

					imported, err := importProps(importedStr, propsMap, &evaluatedImportsBranch)
					if err != nil {
						return nil, err
					}
					for k, v := range imported {
						evaluated[k] = v
					}
				} else {
					return nil, fmt.Errorf("config import statement has invalid value: `%v`", importedInterface)
				}
			}
		} else {
			return nil, fmt.Errorf("config import statement has invalid value: `%v`", value)
		}
		for key, value := range c {
			if key != extendConfigKey {
				evaluated[key] = value
			}
		}
		return evaluated, nil
	}
	return c, nil
}

func importProps(importedStr string, propsMap map[string]configProps, evaluatedImports *[]string) (configProps, error) {
	if contains(*evaluatedImports, importedStr) {
		return nil, fmt.Errorf("circular import in configs (config imports chain: `%s` -> `%s`)", strings.Join(*evaluatedImports, " -> "), importedStr)
	}
	*evaluatedImports = append(*evaluatedImports, importedStr)
	if imported, ok := propsMap[importedStr]; ok {
		return imported.evaluateConfigImports(propsMap, evaluatedImports)
	}
	return nil, fmt.Errorf("trying to import `%s`, but such config does not exist", importedStr)
}

func contains(slice []string, element string) bool {
	for _, e := range slice {
		if element == e {
			return true
		}
	}
	return false
}
