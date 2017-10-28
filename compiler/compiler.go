package compiler

import (
	"regexp"
	"strconv"

	. "github.com/dankraw/ssh-aliases/domain"
)

type Compiler struct {
	expander     *Expander
	groupsRegexp *regexp.Regexp
}

func NewCompiler() *Compiler {
	return &Compiler{
		expander:     NewExpander(),
		groupsRegexp: regexp.MustCompile("%(\\d+)"),
	}
}

type TemplateReplacemenet struct {
	beginIdx       int
	endIdx         int
	replacementIdx int
}

func (c *Compiler) Compile(input HostConfigInput) ([]HostConfigResult, error) {
	results := []HostConfigResult{}
	expanded, err := c.expander.expand(input.HostnamePattern)
	if err != nil {
		return nil, err
	}
	templateGroups := c.groupsRegexp.FindAllStringSubmatchIndex(input.AliasTemplate, -1)
	selectors := []TemplateReplacemenet{}
	for _, group := range templateGroups {
		hostnameGroupSelect, _ := strconv.Atoi(input.AliasTemplate[group[2]:group[3]])
		selectors = append(selectors, TemplateReplacemenet{group[0], group[1], hostnameGroupSelect - 1})
	}
	for _, h := range expanded {
		results = append(results, HostConfigResult{
			Host:       c.compileToTargetHost(input.AliasTemplate, selectors, h),
			HostName:   h.Hostname,
			HostConfig: input.HostConfig,
		})
	}
	return results, nil
}

func (c *Compiler) compileToTargetHost(aliasTemplate string, selectors []TemplateReplacemenet, host ExpandedHostname) string {
	alias := ""
	for i, s := range selectors {
		if i == 0 {
			alias += aliasTemplate[0:s.beginIdx]
		}

		alias += host.Replacements[s.replacementIdx]
		nextIdx := i + 1
		if nextIdx < len(selectors) {
			nextSelector := selectors[nextIdx]
			alias += aliasTemplate[s.endIdx:nextSelector.beginIdx]
		} else {
			alias += aliasTemplate[s.endIdx:]
		}
	}
	return alias
}
