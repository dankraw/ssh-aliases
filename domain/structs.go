package domain

type HostConfigInput struct {
	AliasName       string
	HostnamePattern string
	AliasTemplate   string
	HostConfig      HostConfig
}

type HostConfig map[string]interface{}

func (c HostConfig) ToHostConfigEntries() []HostConfigEntry {
	entries := []HostConfigEntry{}
	for k, v := range c {
		entries = append(entries, HostConfigEntry{k, v})
	}
	return entries
}

type HostConfigResult struct {
	Host       string
	HostName   string
	HostConfig HostConfig
}

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
