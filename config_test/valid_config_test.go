package config

import (
	"testing"

	"github.com/dankraw/ssh-aliases/compiler"
	"github.com/dankraw/ssh-aliases/config"
	"github.com/stretchr/testify/assert"
)

func TestShouldReadCompleteConfigFromDir(t *testing.T) {
	t.Parallel()

	// given
	reader := config.NewReader()

	// when
	ctx, err := reader.ReadConfigs("./test_fixtures/valid")

	// then
	assert.NoError(t, err)
	assert.Equal(t, compiler.InputContext{
		Sources: []compiler.ContextSource{
			{
				SourceName: "test_fixtures/valid/example.hcl",
				Hosts: []compiler.ExpandingHostConfig{{
					AliasName:       "service-a",
					HostnamePattern: "service-a[1..5].my.domain1.example.com",
					AliasTemplate:   "a{#1}",
					Config: compiler.ConfigProperties{{
						Key:   "IdentityFile",
						Value: "a_1001_id_secret_rsa.pem",
					}, {
						Key:   "Port",
						Value: 22,
					}, {
						Key:   "User",
						Value: "deployment",
					}},
				}, {
					AliasName:       "service-b",
					HostnamePattern: "service-b[1..2].example.com",
					AliasTemplate:   "b{#1}",
					Config: compiler.ConfigProperties{{
						Key:   "IdentityFile",
						Value: "b_id_1001_rsa.pem",
					}, {
						Key:   "Port",
						Value: 22,
					}},
				}},
			}, {
				SourceName: "test_fixtures/valid/values.hcl",
				Hosts:      []compiler.ExpandingHostConfig{},
			},
		},
	}, ctx)
}
