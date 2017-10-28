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
	config, err := reader.ReadConfigs(FIXTURE_DIR)

	// then
	assert.NoError(t, err)
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
