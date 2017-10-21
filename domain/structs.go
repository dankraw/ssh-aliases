package domain

import "regexp"

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
