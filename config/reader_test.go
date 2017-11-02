package config

import (
	"testing"

	. "github.com/dankraw/ssh-aliases/domain"
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
			SSHConfig: HostConfig{
				"identity_file": "b_id_rsa.pub",
				"port":          22,
			},
		}}, RawSSHConfigs: RawSSHConfigs{
			"service-a": []map[string]interface{}{{
				"identity_file": "a_id_rsa.pub",
				"port":          22,
			}}},
	}, config)
}
