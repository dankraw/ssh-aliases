package compiler

type HostConfigInput struct {
	AliasName       string
	HostnamePattern string
	AliasTemplate   string
	HostConfig      HostConfigEntries
}

type HostConfigResult struct {
	Host       string
	HostName   string
	HostConfig HostConfigEntries
}

type HostConfigEntries []HostConfigEntry

type HostConfigEntry struct {
	Key   string
	Value interface{}
}

type ByHostConfigEntryKey []HostConfigEntry

func (s ByHostConfigEntryKey) Len() int {
	return len(s)
}

func (s ByHostConfigEntryKey) Less(i, j int) bool {
	return s[i].Key < s[j].Key
}

func (s ByHostConfigEntryKey) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
