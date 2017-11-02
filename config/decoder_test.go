package config

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fixtureDir = "./test-fixtures"

func TestShouldDecodeConfig(t *testing.T) {
	t.Parallel()

	// given
	data, _ := ioutil.ReadFile(filepath.Join(fixtureDir, "example.hcl"))

	// when
	config, _ := NewDecoder().Decode(data)

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
			SSHConfig: HostConfig{
				"identity_file": "b_id_rsa.pub",
				"port":          22,
			},
		}}, RawSSHConfigs: RawSSHConfigs{
			"service-a": []map[string]interface{}{{
				"identity_file": "a_id_rsa.pub",
				"port":          22,
			}},
		},
	}, config)

}
