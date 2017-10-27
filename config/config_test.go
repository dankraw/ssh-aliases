package config

import (
	. "github.com/dankraw/ssh-aliases/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldMapToHostConfigInputs(t *testing.T) {
	t.Parallel()

	// given
	config := Config{
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
	}

	// when
	inputs, err := config.ToHostConfigInputs()

	// then
	assert.NoError(t, err)
	assert.Equal(t, []HostConfigInput{{
		HostnamePattern: "service-a[1..5].example.com",
		AliasTemplate:   "a%1",
		HostConfig: HostConfig{
			IdentityFile: "a_id_rsa.pub",
			Port:         22,
		}}, {
		HostnamePattern: "service-b[1..2].example.com",
		AliasTemplate:   "b%1",
		HostConfig: HostConfig{
			IdentityFile: "b_id_rsa.pub",
			Port:         22,
		}},
	}, inputs)
}

func TestShouldReturnErrorOnNotFoundSSHConfig(t *testing.T) {
	t.Parallel()

	// given
	config := Config{
		Aliases: []Alias{{
			Name:          "service-a",
			Pattern:       "service-a[1..5].example.com",
			Template:      "a%1",
			SSHConfigName: "this-does-not-exists",
		}},
	}

	// when
	results, err := config.ToHostConfigInputs()

	// then
	assert.Nil(t, results)
	assert.Error(t, err)
	assert.Equal(t, "No ssh-config named this-does-not-exists found (used by service-a alias)", err.Error())
}
