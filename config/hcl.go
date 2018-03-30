package config

type rawFileContext struct {
	Hosts      []host                 `hcl:"host"`
	RawConfigs map[string]rawConfig   `hcl:"config"`
	Variables  map[string]interface{} `hcl:"var"`
}

type rawConfig []map[string]interface{}

type host struct {
	Name           string      `hcl:",key"`
	Hostname       string      `hcl:"hostname"`
	Alias          string      `hcl:"alias"`
	RawConfigOrRef interface{} `hcl:"config"`
}
