package config

import (
	"testing"

	"github.com/dankraw/ssh-aliases/compiler"
	"github.com/stretchr/testify/assert"
)

func TestShouldMapToHostConfigInputs(t *testing.T) {
	t.Parallel()

	// given
	config := rawConfigContext{
		Hosts: []host{{
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
		}}, RawConfigs: rawConfigs{
			"service-a": rawConfig{{
				"identity_file": "a_id_rsa.pub",
				"port":          22,
			}}},
	}

	// when
	inputs, err := config.toExpandingHostConfigs()

	// then
	assert.NoError(t, err)
	assert.Equal(t, []compiler.ExpandingHostConfig{{
		AliasName:       "service-a",
		HostnamePattern: "service-a[1..5].example.com",
		AliasTemplate:   "a%1",
		Config: compiler.ConfigProperties{
			{Key: "IdentityFile", Value: "a_id_rsa.pub"},
			{Key: "Port", Value: 22},
		}}, {
		AliasName:       "service-b",
		HostnamePattern: "service-b[1..2].example.com",
		AliasTemplate:   "b%1",
		Config: compiler.ConfigProperties{
			{Key: "IdentityFile", Value: "b_id_rsa.pub"},
			{Key: "Port", Value: 22},
		}},
	}, inputs)
}

func TestShouldReturnErrorOnDuplicateKey(t *testing.T) {
	t.Parallel()

	// given
	config := rawConfigContext{
		Hosts: []host{{
			Name:     "service-b",
			Hostname: "service-b[1..2].example.com",
			Alias:    "b%1",
			RawConfigOrRef: []map[string]interface{}{{
				"identity_file": "b_id_rsa.pub",
			}, {
				"identity_file": "c_id_rsa.pub",
			}},
		}}}

	// when
	inputs, err := config.toExpandingHostConfigs()

	// then
	assert.Error(t, err)
	assert.Nil(t, inputs)
	assert.Equal(t, "duplicate config property `identity_file` for host `service-b`", err.Error())
}

func TestShouldReturnErrorOnDuplicateKeyInRawConfigs(t *testing.T) {
	t.Parallel()

	// given
	config := rawConfigContext{
		RawConfigs: rawConfigs{
			"service-a": rawConfig{
				{"identity_file": "abc"},
				{"identity_file": "abc"},
			},
		},
	}

	// when
	inputs, err := config.toExpandingHostConfigs()

	// then
	assert.Error(t, err)
	assert.Nil(t, inputs)
	assert.Equal(t, "duplicate config entry `identity_file` in host `service-a`", err.Error())
}

func TestShouldReturnErrorOnNotFoundSSHConfig(t *testing.T) {
	t.Parallel()

	// given
	config := rawConfigContext{
		Hosts: []host{{
			Name:           "service-a",
			Hostname:       "service-a[1..5].example.com",
			Alias:          "a%1",
			RawConfigOrRef: "this-does-not-exists",
		}},
	}

	// when
	results, err := config.toExpandingHostConfigs()

	// then
	assert.Nil(t, results)
	assert.Error(t, err)
	assert.Equal(t, "no config `this-does-not-exists` found (used by host `service-a`)", err.Error())
}

func TestShouldMergeWithOtherConfig(t *testing.T) {
	t.Parallel()

	// given
	config := rawConfigContext{
		Hosts: []host{{
			Name: "project1",
		}},
		RawConfigs: rawConfigs{
			"project1-config": rawConfig{{
				"identity_file": "a_id_rsa.pub",
			}},
		},
	}

	// when
	merged, err := mergeRawConfigContexts(config, rawConfigContext{
		Hosts: []host{{
			Name: "project2",
		}},
		RawConfigs: rawConfigs{
			"project2-config": rawConfig{{
				"port": 22,
			}},
		},
	})

	// then
	assert.NoError(t, err)
	assert.Equal(t, rawConfigContext{
		Hosts: []host{{
			Name: "project1",
		}, {
			Name: "project2",
		}},
		RawConfigs: rawConfigs{
			"project1-config": rawConfig{{
				"identity_file": "a_id_rsa.pub",
			}},
			"project2-config": rawConfig{{
				"port": 22,
			}},
		},
	}, merged)
}

func TestShouldReturnErrorOnDuplicateSSHConfigWhenMerging(t *testing.T) {
	t.Parallel()

	// given
	config := rawConfigContext{
		RawConfigs: rawConfigs{
			"service-a": rawConfig{{
				"identity_file": "a_id_rsa.pub",
			}},
		},
	}
	config2 := rawConfigContext{
		RawConfigs: rawConfigs{
			"service-a": rawConfig{{
				"port": 22,
			}},
		},
	}

	// when
	_, err := mergeRawConfigContexts(config, config2)

	// then
	assert.Error(t, err)
	assert.Equal(t, "duplicate config `service-a`", err.Error())
}

func TestShouldReturnErrorOnDuplicateAlias(t *testing.T) {
	t.Parallel()

	// given
	config := rawConfigContext{
		Hosts: []host{{
			Name: "service-a",
		}, {
			Name: "service-a",
		}},
	}

	// when
	results, err := config.toExpandingHostConfigs()

	// then
	assert.Nil(t, results)
	assert.Error(t, err)
}

func TestShouldSortHostConfigAndSanitizeKeywords(t *testing.T) {
	t.Parallel()

	// given
	c := configProps{
		"b": "something",
		"c": "abc",
		"d": 0,
		"a": 123,
	}

	// when
	entries := c.toSortedProperties()

	// then
	assert.Equal(t, compiler.ConfigProperties{
		{Key: "A", Value: 123},
		{Key: "B", Value: "something"},
		{Key: "C", Value: "abc"},
		{Key: "D", Value: 0},
	}, entries)
}
