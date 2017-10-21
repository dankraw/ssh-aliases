package compiler

import (
	"fmt"
	"strings"
	. "github.com/dankraw/ssh-aliases/domain"
)

type Compiler struct {}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(input HostConfigInput) ([]HostConfigResult, error) {
	results := []HostConfigResult{}
	for _, host := range input.Hostnames {
		results = append(results, HostConfigResult{
			Host: c.compileToTargetHost(input, host),
			HostConfig: input.HostConfig,
		})
	}
	return results, nil
}

func (c *Compiler) compileToTargetHost(input HostConfigInput, host string) string {
	found := input.HostnameRegexp.FindStringSubmatch(host)
	target := input.TargetPatternTemplate
	for i, group := range found {
		groupPlaceholder := fmt.Sprintf("%%%d", i)
		target = strings.Replace(target, groupPlaceholder, group, -1)
	}
	return target
}
