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
	ctx, err := reader.ReadConfigs("./test_fixtures/valid/basic_with_variables")

	// then
	assert.NoError(t, err)
	assert.Equal(t, compiler.InputContext{
		Sources: []compiler.ContextSource{
			{
				SourceName: "test_fixtures/valid/basic_with_variables/example.hcl",
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
				SourceName: "test_fixtures/valid/basic_with_variables/variables.hcl",
				Hosts:      []compiler.ExpandingHostConfig{},
			},
		},
	}, ctx)
}

func TestShouldReadFilesWithImportedConfigs(t *testing.T) {
	t.Parallel()

	// given
	reader := config.NewReader()

	// when
	ctx, err := reader.ReadConfigs("./test_fixtures/valid/importing_configs")

	// then
	assert.NoError(t, err)
	assert.Equal(t, compiler.InputContext{
		Sources: []compiler.ContextSource{
			{
				SourceName: "test_fixtures/valid/importing_configs/example.hcl",
				Hosts: []compiler.ExpandingHostConfig{{
					AliasName:       "abc",
					HostnamePattern: "servcice-abc.example.com",
					AliasTemplate:   "abc",
					Config: compiler.ConfigProperties{{
						Key:   "Additional",
						Value: "extension",
					}, {
						Key:   "Another",
						Value: "one",
					}, {
						Key:   "X",
						Value: "y",
					}},
				}, {
					AliasName:       "def",
					HostnamePattern: "servcice-def.example.com",
					AliasTemplate:   "def",
					Config: compiler.ConfigProperties{{
						Key:   "Additional",
						Value: "extension",
					}, {
						Key:   "Another",
						Value: "one",
					}, {
						Key:   "SomeProp",
						Value: 123,
					}, {
						Key:   "This",
						Value: "happens",
					}},
				}},
			},
		},
	}, ctx)
}

func TestShouldReadHostDefinitionsWithoutHostnames(t *testing.T) {
	t.Parallel()

	// given
	reader := config.NewReader()

	// when
	ctx, err := reader.ReadConfigs("./test_fixtures/valid/no_hostname")

	// then
	assert.NoError(t, err)
	assert.Equal(t, compiler.InputContext{
		Sources: []compiler.ContextSource{
			{
				SourceName: "test_fixtures/valid/no_hostname/wildcard.hcl",
				Hosts: []compiler.ExpandingHostConfig{{
					AliasName:     "all",
					AliasTemplate: "*",
					Config: compiler.ConfigProperties{{
						Key:   "A",
						Value: 1,
					}},
				}},
			},
		},
	}, ctx)
}
