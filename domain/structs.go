package domain

type HostConfigInput struct {
	HostnamePattern string
	AliasTemplate   string
	HostConfig      HostConfig
}

type HostConfig struct {
	IdentityFile string
	Port         int
}

type HostConfigResult struct {
	Host       string
	HostName   string
	HostConfig HostConfig
}
