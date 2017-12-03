package compiler

// ExpandingHostConfig is the input for the ssh-aliases compiler
type ExpandingHostConfig struct {
	AliasName       string
	HostnamePattern string
	AliasTemplate   string
	Config          ConfigProperties
}

// InputContext is the container for all host and configs
type InputContext struct {
	Sources []ContextSource
}

// ContextSource is represents a single piece of source that provides host and configs definitions
type ContextSource struct {
	SourceName string
	Hosts      []ExpandingHostConfig
}

// HostEntity is the outcome of ssh-alises compiler
type HostEntity struct {
	Host     string
	HostName string
	Config   ConfigProperties
}

// ConfigProperties is a list of ssh config properties
type ConfigProperties []ConfigProperty

// ConfigProperty is a key-value container
type ConfigProperty struct {
	Key   string
	Value interface{}
}

// ByConfigPropertyKey can be used to sort an array of ConfigProperties by their keys
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
