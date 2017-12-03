package config

import (
	"testing"

	"github.com/dankraw/ssh-aliases/compiler"
	"github.com/stretchr/testify/assert"
)

func TestShouldReadCompleteConfigFromDir(t *testing.T) {
	t.Parallel()

	// given
	reader := NewReader()

	// when
	ctx, err := reader.ReadConfigs(fixtureDir)

	// then
	assert.NoError(t, err)
	assert.Equal(t, compiler.InputContext{
		Sources: []compiler.ContextSource{
			{
				SourceName: "test-fixtures/example.hcl",
				Hosts: []compiler.ExpandingHostConfig{{
					AliasName:       "service-a",
					HostnamePattern: "service-a[1..5].example.com",
					AliasTemplate:   "a{#1}",
					Config: compiler.ConfigProperties{{
						Key:   "IdentityFile",
						Value: "a_id_rsa.pem",
					}, {
						Key:   "Port",
						Value: 22,
					}},
				}, {
					AliasName:       "service-b",
					HostnamePattern: "service-b[1..2].example.com",
					AliasTemplate:   "b{#1}",
					Config: compiler.ConfigProperties{{
						Key:   "IdentityFile",
						Value: "b_id_rsa.pem",
					}, {
						Key:   "Port",
						Value: 22,
					}},
				}},
			},
		},
	}, ctx)
}
