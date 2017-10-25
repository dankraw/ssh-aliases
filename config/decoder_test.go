package config

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fixtureDir = "./test-fixtures"

func TestShouldDecodeConfig(t *testing.T) {
	// given
	data, _ := ioutil.ReadFile(filepath.Join(fixtureDir, "example.hcl"))

	// when
	config, _ := NewDecoder().decode(data)

	// then
	assert.Equal(t, Config{
		Aliases: []Alias{{
			Name:          "service-a",
			Pattern:       "service-a[1..5].example.com",
			Template:      "a%1",
			SSHConfigName: "service-a",
		}, {
			Name:     "service-b",
			Pattern:  "service-b[1..2].example.com",
			Template: "b%1",
			SSHConfig: EmbeddedSSHConfig{
				IdentityFile: "b_id_rsa.pub",
				Port:         22,
			},
		}}, SSHConfigs: []SSHConfig{{
			Name:         "service-a",
			IdentityFile: "a_id_rsa.pub",
			Port:         22,
		}},
	}, config)
}
