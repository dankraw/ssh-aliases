package config

import (
	"testing"

	"github.com/dankraw/ssh-aliases/compiler"
	"github.com/stretchr/testify/assert"
)

func TestShouldConvertToCompilerInputContext(t *testing.T) {
	t.Parallel()

	// given
	ctx := rawDirContext{
		RawSources: []rawContextSource{{
			SourceName: "a",
			RawContext: rawFileContext{
				Hosts: []host{{
					Name:           "service-a",
					Hostname:       "service-a[1..5].example.com",
					Alias:          "a{#1}",
					RawConfigOrRef: "service-a",
				}, {
					Name:     "service-b",
					Hostname: "service-b[1..2].example.com",
					Alias:    "b{#1}",
					RawConfigOrRef: []map[string]interface{}{{
						"identity_file": "b_id_rsa.pub",
					}, {
						"port": 22,
					}},
				}},
				RawConfigs: rawConfigs{
					"service-a": rawConfig{{
						"identity_file": "a_id_rsa.pub",
						"port":          22,
					}},
				},
			},
		}, {
			SourceName: "b",
			RawContext: rawFileContext{
				Hosts: []host{{
					Name:     "service-b",
					Hostname: "service-b.example.com",
					Alias:    "b",
					RawConfigOrRef: []map[string]interface{}{{
						"identity_file": "b_id_rsa.pub",
					}, {
						"port": 22,
					}},
				}},
			},
		}},
	}

	// when
	compilerCtx, err := ctx.toCompilerInputContext()

	// then
	assert.NoError(t, err)
	assert.Equal(t, compiler.InputContext{
		Sources: []compiler.ContextSource{{
			SourceName: "a",
			Hosts: []compiler.ExpandingHostConfig{{
				AliasName:       "service-a",
				HostnamePattern: "service-a[1..5].example.com",
				AliasTemplate:   "a{#1}",
				Config: compiler.ConfigProperties{
					{Key: "IdentityFile", Value: "a_id_rsa.pub"},
					{Key: "Port", Value: 22},
				},
			}, {
				AliasName:       "service-b",
				HostnamePattern: "service-b[1..2].example.com",
				AliasTemplate:   "b{#1}",
				Config: compiler.ConfigProperties{
					{Key: "IdentityFile", Value: "b_id_rsa.pub"},
					{Key: "Port", Value: 22},
				},
			}},
		}, {
			SourceName: "b",
			Hosts: []compiler.ExpandingHostConfig{{
				AliasName:       "service-b",
				HostnamePattern: "service-b.example.com",
				AliasTemplate:   "b",
				Config: compiler.ConfigProperties{
					{Key: "IdentityFile", Value: "b_id_rsa.pub"},
					{Key: "Port", Value: 22},
				},
			}},
		}},
	}, compilerCtx)
}

func TestShouldReturnErrorOnDuplicateKeyInEmbeddedConfig(t *testing.T) {
	t.Parallel()

	// given
	ctx := rawDirContext{
		RawSources: []rawContextSource{{
			SourceName: "b",
			RawContext: rawFileContext{
				Hosts: []host{{
					Name:     "service-b",
					Hostname: "service-b[1..2].example.com",
					Alias:    "b{#1}",
					RawConfigOrRef: []map[string]interface{}{{
						"identity_file": "b_id_rsa.pub",
					}, {
						"identity_file": "c_id_rsa.pub",
					}},
				}},
			},
		}},
	}

	// when
	compilerCtx, err := ctx.toCompilerInputContext()

	// then
	assert.Error(t, err)
	assert.Equal(t, compiler.InputContext{}, compilerCtx)
	assert.Equal(t, "duplicate config property `identity_file` for host `service-b`", err.Error())
}

func TestShouldReturnErrorOnDuplicateKeyInRawConfigs(t *testing.T) {
	t.Parallel()

	// given
	ctx := rawDirContext{
		RawSources: []rawContextSource{{
			SourceName: "a",
			RawContext: rawFileContext{
				RawConfigs: rawConfigs{
					"service-a": rawConfig{
						{"identity_file": "abc"},
						{"identity_file": "abc"},
					},
				},
			},
		}},
	}

	// when
	compilerCtx, err := ctx.toCompilerInputContext()

	// then
	assert.Error(t, err)
	assert.Equal(t, compiler.InputContext{}, compilerCtx)
	assert.Equal(t, "duplicate config entry `identity_file` in host `service-a`", err.Error())
}

func TestShouldReturnErrorOnNotFoundSSHConfig(t *testing.T) {
	t.Parallel()

	// given
	ctx := rawDirContext{
		RawSources: []rawContextSource{{
			SourceName: "a",
			RawContext: rawFileContext{
				Hosts: []host{{
					Name:           "service-a",
					Hostname:       "service-a[1..5].example.com",
					Alias:          "a{#1}",
					RawConfigOrRef: "this-does-not-exists",
				}},
			},
		}},
	}

	// when
	compilerCtx, err := ctx.toCompilerInputContext()

	// then
	assert.Error(t, err)
	assert.Equal(t, compiler.InputContext{}, compilerCtx)
	assert.Equal(t, "no config `this-does-not-exists` found (used by host `service-a`)", err.Error())
}
func TestShouldReturnErrorOnDuplicateSSHConfig(t *testing.T) {
	t.Parallel()

	// given
	ctx := rawDirContext{
		RawSources: []rawContextSource{{
			SourceName: "a",
			RawContext: rawFileContext{
				RawConfigs: rawConfigs{
					"config-a": rawConfig{{
						"identity_file": "a_id_rsa.pub",
					}},
				},
			},
		}, {
			SourceName: "b",
			RawContext: rawFileContext{
				RawConfigs: rawConfigs{
					"config-a": rawConfig{{
						"port": 22,
					}},
				},
			},
		}},
	}

	// when
	compilerCtx, err := ctx.toCompilerInputContext()

	// then
	assert.Error(t, err)
	assert.Equal(t, compiler.InputContext{}, compilerCtx)
	assert.Equal(t, "duplicate config with name `config-a`", err.Error())
}

func TestShouldReturnErrorOnDuplicateAlias(t *testing.T) {
	t.Parallel()

	// given
	ctx := rawDirContext{
		RawSources: []rawContextSource{{
			SourceName: "a",
			RawContext: rawFileContext{
				Hosts: []host{{
					Name:           "service-a",
					RawConfigOrRef: "config-a",
				}, {
					Name:           "service-a",
					RawConfigOrRef: "config-a",
				}},
				RawConfigs: rawConfigs{
					"config-a": rawConfig{{
						"port": 22,
					}},
				},
			},
		}},
	}

	// when
	compilerCtx, err := ctx.toCompilerInputContext()

	// then
	assert.Equal(t, compiler.InputContext{}, compilerCtx)
	assert.Error(t, err)
	assert.Equal(t, "duplicate host `service-a`", err.Error())
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
