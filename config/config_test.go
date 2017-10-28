package config

import (
	"testing"

	. "github.com/dankraw/ssh-aliases/domain"
	"github.com/stretchr/testify/assert"
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
		AliasName:       "service-a",
		HostnamePattern: "service-a[1..5].example.com",
		AliasTemplate:   "a%1",
		HostConfig: HostConfig{
			IdentityFile: "a_id_rsa.pub",
			Port:         22,
		}}, {
		AliasName:       "service-b",
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

func TestShouldReturnErrorOnDuplicateSSHConfig(t *testing.T) {
	t.Parallel()

	// given
	config := Config{
		SSHConfigs: []SSHConfig{{
			Name: "service-a",
		}, {
			Name: "service-a",
		}},
	}

	// when
	results, err := config.ToHostConfigInputs()

	// then
	assert.Nil(t, results)
	assert.Error(t, err)
	assert.Equal(t, "Duplicate ssh-config with name service-a", err.Error())
}

func TestShouldReturnErrorOnDuplicateAlias(t *testing.T) {
	t.Parallel()

	// given
	config := Config{
		Aliases: []Alias{{
			Name: "project1",
		}},
		SSHConfigs: []SSHConfig{{
			Name: "service-a",
		}},
	}

	// when
	config.Merge(Config{
		Aliases: []Alias{{
			Name: "project2",
		}},
		SSHConfigs: []SSHConfig{{
			Name: "service-b",
		}},
	})

	// then
	assert.Equal(t, Config{
		Aliases: []Alias{{
			Name: "project1",
		}, {
			Name: "project2",
		}},
		SSHConfigs: []SSHConfig{{
			Name: "service-a",
		}, {
			Name: "service-b",
		}},
	}, config)
}

func TestShouldMergeWithOtherConfig(t *testing.T) {
	t.Parallel()

	// given
	config := Config{
		Aliases: []Alias{{
			Name: "service-a",
		}, {
			Name: "service-a",
		}},
	}

	// when
	results, err := config.ToHostConfigInputs()

	// then
	assert.Nil(t, results)
	assert.Error(t, err)
}
