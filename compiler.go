package main

import (
	"fmt"
	"regexp"
	"strings"
)

type HostConfigInput struct {
	Hostnames []string
	HostnameRegexp *regexp.Regexp
	TargetPatternTemplate string
	HostConfig *HostConfig
}

type HostConfig struct {
	IdentityFile string
	Port uint16
}

type HostConfigResult struct {
	Host string
	HostConfig *HostConfig
}

func compile(input HostConfigInput) ([]HostConfigResult, error) {
	results := []HostConfigResult{}
	for _, host := range input.Hostnames {
		results = append(results, HostConfigResult{
			Host: compileToTargetHost(input, host),
			HostConfig: input.HostConfig,
		})
	}
	return results, nil
}

func compileToTargetHost(input HostConfigInput, host string) string {
	found := input.HostnameRegexp.FindStringSubmatch(host)
	target := input.TargetPatternTemplate
	for i, group := range found {
		groupPlaceholder := fmt.Sprintf("%%%d", i)
		target = strings.Replace(target, groupPlaceholder, group, -1)
	}
	return target
}
