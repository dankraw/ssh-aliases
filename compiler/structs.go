package compiler

type ExpandingHostConfig struct {
	AliasName       string
	HostnamePattern string
	AliasTemplate   string
	Config          ConfigProperties
}

type HostEntity struct {
	Host     string
	HostName string
	Config   ConfigProperties
}

type ConfigProperties []ConfigProperty

type ConfigProperty struct {
	Key   string
	Value interface{}
}

type ByConfigPropertyKey []ConfigProperty

func (s ByConfigPropertyKey) Len() int {
	return len(s)
}

func (s ByConfigPropertyKey) Less(i, j int) bool {
	return s[i].Key < s[j].Key
}

func (s ByConfigPropertyKey) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
