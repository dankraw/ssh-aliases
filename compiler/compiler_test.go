package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	t.Parallel()

	// given
	sshConfig := ConfigProperties{{"identity_file", "~/.ssh/id_rsa"}}
	input := ExpandingHostConfig{
		HostnamePattern: "x-master[1..2].myproj-prod.dc1.net",
		AliasTemplate:   "host{#1}-dc1",
		Config:          sshConfig,
	}

	// when
	results, err := NewCompiler().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []HostEntity{{
		Host:     "host1-dc1",
		HostName: "x-master1.myproj-prod.dc1.net",
		Config:   sshConfig,
	}, {
		Host:     "host2-dc1",
		HostName: "x-master2.myproj-prod.dc1.net",
		Config:   sshConfig,
	}}, results)
}

func TestShouldReplaceAllGroupMatchOccurrences(t *testing.T) {
	t.Parallel()

	// given
	input := ExpandingHostConfig{
		HostnamePattern: "x-[master1].myproj-prod.dc1.net",
		AliasTemplate:   "{#1}-{#1}-{#1}",
	}

	// when
	results, err := NewCompiler().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "master1-master1-master1", results[0].Host)
}

func TestShouldExpandHostnameWithProvidedRange(t *testing.T) {
	t.Parallel()

	// given
	input := ExpandingHostConfig{
		HostnamePattern: "x-master[4..6].myproj-prod.dc1.net",
		AliasTemplate:   "m{#1}",
	}

	// when
	results, err := NewCompiler().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []HostEntity{{
		Host:     "m4",
		HostName: "x-master4.myproj-prod.dc1.net",
	}, {
		Host:     "m5",
		HostName: "x-master5.myproj-prod.dc1.net",
	}, {
		Host:     "m6",
		HostName: "x-master6.myproj-prod.dc1.net",
	}}, results)
}

func TestShouldAllowStaticAliasDefinitions(t *testing.T) {
	t.Parallel()

	// given
	input := ExpandingHostConfig{
		HostnamePattern: "x-master1.myproj-prod.dc1.net",
		AliasTemplate:   "master1",
	}

	// when
	results, err := NewCompiler().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "master1", results[0].Host)
}

func TestShouldAllowHostDefinitionsWithoutHostnames(t *testing.T) {
	t.Parallel()

	// given
	input := ExpandingHostConfig{
		AliasTemplate: "*",
	}

	// when
	results, err := NewCompiler().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "*", results[0].Host)
	assert.Equal(t, "", results[0].HostName)
}

func TestRegexpCompile(t *testing.T) {
	t.Parallel()

	// given
	sshConfig := ConfigProperties{{"identity_file", "~/.ssh/id_rsa"}}
	input := ExpandingHostConfig{
		HostnamePattern: "x-master(\\d+)\\.myproj-([a-z]+)\\.dc1\\.net",
		AliasTemplate:   "{#2}.host{#1}.dc1",
		Config:          sshConfig,
	}
	hosts := InputHosts{
		"y-master1.myproj-prod.dc2",
		"x-master2.myproj-prod-dc1.net",
		"x-master3.myproj-prod.dc1.net",
		"x-master4.other-prod.dc1.net",
		"x-master5.myproj-test.dc1.net",
		"x-master6.myproj-test.dc1.net x-master7.myproj-dev.dc1.net ddd",
	}

	// when
	results, err := NewCompiler().CompileRegexp(input, hosts)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []HostEntity{{
		Host:     "prod.host3.dc1",
		HostName: "x-master3.myproj-prod.dc1.net",
		Config:   sshConfig,
	}, {
		Host:     "test.host5.dc1",
		HostName: "x-master5.myproj-test.dc1.net",
		Config:   sshConfig,
	}, {
		Host:     "test.host6.dc1",
		HostName: "x-master6.myproj-test.dc1.net",
		Config:   sshConfig,
	}, {
		Host:     "dev.host7.dc1",
		HostName: "x-master7.myproj-dev.dc1.net",
		Config:   sshConfig,
	}}, results)
}
