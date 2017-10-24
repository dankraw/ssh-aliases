package compiler

import (
	"fmt"
	"strings"
	. "github.com/dankraw/ssh-aliases/domain"
)

type Compiler struct {
	expander *Expander
}

func NewCompiler() *Compiler {
	return &Compiler{
		NewExpander(),
	}
}

func (c *Compiler) Compile(input HostConfigInput) ([]HostConfigResult, error) {
	results := []HostConfigResult{}
	for _, host := range input.Hostnames {
		expanded, err := c.expander.expand(host)
		if err != nil {
			return nil, err
		}
		for _, h := range expanded {
			results = append(results, HostConfigResult{
				Host: c.compileToTargetHost(input, h),
				HostName: h,
				HostConfig: input.HostConfig,
			})
		}
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

