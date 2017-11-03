package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReadCompleteConfigFromDir(t *testing.T) {
	t.Parallel()

	// given
	reader := NewReader()

	// when
	config, err := reader.ReadConfigs(fixtureDir)

	// then
	assert.NoError(t, err)
	assert.Equal(t, HostsWithConfigs{
		Hosts: []Host{{
			Name:        "service-a",
			Hostname:    "service-a[1..5].example.com",
			Alias:       "a%1",
			ConfigOrRef: "service-a",
		}, {
			Name:     "service-b",
			Hostname: "service-b[1..2].example.com",
			Alias:    "b%1",
			ConfigOrRef: []map[string]interface{}{{
				"identity_file": "b_id_rsa.pub",
			}, {
				"port": 22,
			}},
		}}, RawSSHConfigs: RawSSHConfigs{
			"service-a": []map[string]interface{}{{
				"identity_file": "a_id_rsa.pub",
				"port":          22,
			}}},
	}, config)
}
