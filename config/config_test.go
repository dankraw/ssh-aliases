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
			SSHConfig: HostConfig{
				"identity_file": "b_id_rsa.pub",
				"port":          22,
			},
		}}, RawSSHConfigs: RawSSHConfigs{
			"service-a": []map[string]interface{}{{
				"identity_file": "a_id_rsa.pub",
				"port":          22,
			}}},
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
			"identity_file": "a_id_rsa.pub",
			"port":          22,
		}}, {
		AliasName:       "service-b",
		HostnamePattern: "service-b[1..2].example.com",
		AliasTemplate:   "b%1",
		HostConfig: HostConfig{
			"identity_file": "b_id_rsa.pub",
			"port":          22,
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

func TestShouldMergeWithOtherConfig(t *testing.T) {
	t.Parallel()

	// given
	config := Config{
		Aliases: []Alias{{
			Name: "project1",
		}},
		RawSSHConfigs: RawSSHConfigs{
			"project1-config": []map[string]interface{}{{
				"identity_file": "a_id_rsa.pub",
			}},
		},
	}

	// when
	err := config.Merge(Config{
		Aliases: []Alias{{
			Name: "project2",
		}},
		RawSSHConfigs: RawSSHConfigs{
			"project2-config": []map[string]interface{}{{
				"port": 22,
			}},
		},
	})

	// then
	assert.NoError(t, err)
	assert.Equal(t, Config{
		Aliases: []Alias{{
			Name: "project1",
		}, {
			Name: "project2",
		}},
		RawSSHConfigs: RawSSHConfigs{
			"project1-config": []map[string]interface{}{{
				"identity_file": "a_id_rsa.pub",
			}},
			"project2-config": []map[string]interface{}{{
				"port": 22,
			}},
		},
	}, config)
}

func TestShouldReturnErrorOnDuplicateSSHConfigWhenMerging(t *testing.T) {
	t.Parallel()

	// given
	config := Config{
		RawSSHConfigs: RawSSHConfigs{
			"service-a": []map[string]interface{}{{
				"identity_file": "a_id_rsa.pub",
			}},
		},
	}
	config2 := Config{
		RawSSHConfigs: RawSSHConfigs{
			"service-a": []map[string]interface{}{{
				"port": 22,
			}},
		},
	}

	// when
	err := config.Merge(config2)

	// then
	assert.Error(t, err)
	assert.Equal(t, "Duplicate ssh-config with name service-a", err.Error())
}

func TestShouldReturnErrorOnDuplicateAlias(t *testing.T) {
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
