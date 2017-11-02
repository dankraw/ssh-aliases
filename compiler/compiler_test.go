package compiler

import (
	"testing"

	. "github.com/dankraw/ssh-aliases/domain"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	t.Parallel()

	// given
	sshConfig := HostConfigEntries{{"identity_file", "~/.ssh/id_rsa"}}
	input := HostConfigInput{
		HostnamePattern: "x-master[1..2].myproj-prod.dc1.net",
		AliasTemplate:   "host%1-dc1",
		HostConfig:      sshConfig,
	}

	// when
	results, err := NewCompiler().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []HostConfigResult{{
		Host:       "host1-dc1",
		HostName:   "x-master1.myproj-prod.dc1.net",
		HostConfig: sshConfig,
	}, {
		Host:       "host2-dc1",
		HostName:   "x-master2.myproj-prod.dc1.net",
		HostConfig: sshConfig,
	}}, results)
}

func TestShouldReplaceAllGroupMatchOccurrences(t *testing.T) {
	t.Parallel()

	// given
	input := HostConfigInput{
		HostnamePattern: "x-[master1].myproj-prod.dc1.net",
		AliasTemplate:   "%1-%1-%1",
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
	input := HostConfigInput{
		HostnamePattern: "x-master[4..6].myproj-prod.dc1.net",
		AliasTemplate:   "m%1",
	}

	// when
	results, err := NewCompiler().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []HostConfigResult{{
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
	input := HostConfigInput{
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
