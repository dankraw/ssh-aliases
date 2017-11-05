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
			Name:           "service-a",
			Hostname:       "service-a[1..5].example.com",
			Alias:          "a%1",
			RawConfigOrRef: "service-a",
		}, {
			Name:     "service-b",
			Hostname: "service-b[1..2].example.com",
			Alias:    "b%1",
			RawConfigOrRef: []map[string]interface{}{{
				"identity_file": "b_id_rsa.pub",
			}, {
				"port": 22,
			}},
		}}, RawConfigs: RawConfigs{
			"service-a": []map[string]interface{}{{
				"identity_file": "a_id_rsa.pub",
				"port":          22,
			}}},
	}, config)
}
