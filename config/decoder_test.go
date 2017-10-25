package config

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fixtureDir = "./test-fixtures"

func TestShouldLoadConfiguration(t *testing.T) {
	// given
	data, _ := ioutil.ReadFile(filepath.Join(fixtureDir, "example.hcl"))

	// when
	config, _ := NewDecoder().decode(data)

	// then
	assert.Equal(t, Config{
		Aliases: []Alias {{
			Name:          "hermes-frontend",
			Patterns:      []string{"host[1..5].example.com"},
			RegExp:        "(host\\d+)",
			Template:      "%1",
			SSHConfigName: "private",
		}, {
			Name:          "hermes-consumers",
			Patterns:      []string{"host[1..3].example.com"},
			RegExp:        "(host\\d+)",
			Template:      "%1",
			SSHConfigName: "private",
		}}, SSHConfigs: []SSHConfig{{
			Name: "private",
			IdentityFile: "id_rsa.pub",
			Port: 22,
		}},
	}, config)
}