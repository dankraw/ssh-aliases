package compiler

import (
	"fmt"
	"regexp"
	"strconv"
)

// Compiler is responsible for transforming ExpandingHostConfigs into an array of HostEntities.
type Compiler struct {
	expander     *expander
	groupsRegexp *regexp.Regexp
}

// NewCompiler creates an instance of Compiler
func NewCompiler() *Compiler {
	return &Compiler{
		expander:     newExpander(),
		groupsRegexp: regexp.MustCompile("{#(\\d+)}"),
	}
}

type templateReplacement struct {
	beginIdx       int
	endIdx         int
	replacementIdx int
}

// Compile converts a single ExpandingHostConfig into list of HostEntities
func (c *Compiler) Compile(input ExpandingHostConfig) ([]HostEntity, error) {
	if input.HostnamePattern == "" {
		return []HostEntity{{
			Host:   input.AliasTemplate,
			Config: input.Config,
		}}, nil
	}
	expanded, err := c.expander.expand(input.HostnamePattern)
	if err != nil {
		return nil, err
	}
	replacements := c.aliasReplacementGroups(input.AliasTemplate)
	var results []HostEntity
	for _, h := range expanded {
		results = append(results, HostEntity{
			Host:     c.compileToTargetHost(input.AliasTemplate, replacements, h),
			HostName: h.Hostname,
			Config:   input.Config,
		})
	}
	return results, nil
}

func (c *Compiler) compileToTargetHost(aliasTemplate string, replacements []templateReplacement, host expandedHostname) string {
	if len(replacements) == 0 {
		return aliasTemplate
	}
	alias := ""
	for i, s := range replacements {
		if i == 0 {
			alias += aliasTemplate[0:s.beginIdx]
		}

		alias += host.Replacements[s.replacementIdx]
		nextIdx := i + 1
		if nextIdx < len(replacements) {
			nextSelector := replacements[nextIdx]
			alias += aliasTemplate[s.endIdx:nextSelector.beginIdx]
		} else {
			alias += aliasTemplate[s.endIdx:]
		}
	}
	return alias
}

// CompileRegexp compiles regexp ExpandingHostConfig against provided InputHosts
func (c *Compiler) CompileRegexp(input ExpandingHostConfig, hosts InputHosts) ([]HostEntity, error) {
	re, err := regexp.Compile(input.HostnamePattern)
	if err != nil {
		return nil, fmt.Errorf("error compiling hostname pattern of %s: %s", input.AliasName, err.Error())
	}
	replacements := c.aliasReplacementGroups(input.AliasTemplate)
	var results []HostEntity
	for _, host := range hosts {
		match := re.FindAllStringSubmatch(host, -1)
		if match != nil {
			for _, matchedHost := range match {
				h := expandedHostname{
					Hostname:     matchedHost[0],
					Replacements: matchedHost[1:],
				}
				results = append(results, HostEntity{
					Host:     c.compileToTargetHost(input.AliasTemplate, replacements, h),
					HostName: h.Hostname,
					Config:   input.Config,
				})
			}
		}
	}
	return results, nil
}

func (c *Compiler) aliasReplacementGroups(aliasTemplate string) []templateReplacement {
	templateGroups := c.groupsRegexp.FindAllStringSubmatchIndex(aliasTemplate, -1)
	var replacements []templateReplacement
	for _, group := range templateGroups {
		hostnameGroupSelect, _ := strconv.Atoi(aliasTemplate[group[2]:group[3]])
		replacements = append(replacements, templateReplacement{group[0], group[1], hostnameGroupSelect - 1})
	}
	return replacements
}
