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
		groupsRegexp: regexp.MustCompile(`{#(\d+)}`),
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
	if input.AliasTemplate == "" {
		return []HostEntity{{
			Host:   input.HostnamePattern,
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
		alias, err := c.compileToTargetHost(input.AliasTemplate, replacements, h, input.HostnamePattern)
		if err != nil {
			return nil, fmt.Errorf("error compiling host `%s`: %s", input.AliasName, err.Error())
		}
		results = append(results, HostEntity{
			Host:     alias,
			HostName: h.Hostname,
			Config:   input.Config,
		})
	}
	return results, nil
}

func (c *Compiler) compileToTargetHost(aliasTemplate string, replacements []templateReplacement, host expandedHostname, hostnamePattern string) (string, error) {
	if len(replacements) == 0 {
		return aliasTemplate, nil
	}
	err := c.validateReplacements(aliasTemplate, hostnamePattern, replacements, host.Replacements)
	if err != nil {
		return "", err
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
	return alias, nil
}

func (c *Compiler) validateReplacements(aliasTemplate string, hostnamePattern string, aliasReplacements []templateReplacement, patternReplacements []string) error {
	maxIdxAllowed := len(patternReplacements)
	for _, replacement := range aliasReplacements {
		replacementIdx := replacement.replacementIdx + 1
		if replacementIdx > maxIdxAllowed {
			return fmt.Errorf("alias `%s` contains placeholder with index `#%d` being out of bounds, `%s` allows `#%d` as the maximum index",
				aliasTemplate, replacementIdx, hostnamePattern, maxIdxAllowed)
		}
	}
	return nil
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
		for _, matchedHost := range match {
			h := expandedHostname{
				Hostname:     matchedHost[0],
				Replacements: matchedHost[1:],
			}
			alias, err := c.compileToTargetHost(input.AliasTemplate, replacements, h, input.HostnamePattern)
			if err != nil {
				return nil, fmt.Errorf("error compiling regexp host `%s`: %s", input.AliasName, err.Error())
			}
			results = append(results, HostEntity{
				Host:     alias,
				HostName: h.Hostname,
				Config:   input.Config,
			})
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
