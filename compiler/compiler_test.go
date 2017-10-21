package compiler

import (
	"regexp"
	"testing"
	"github.com/stretchr/testify/assert"
	. "github.com/dankraw/ssh-aliases/domain"
)

func TestCompile(t *testing.T) {
	t.Parallel()

	// given
	r, _ := regexp.Compile("x-([a-z]+\\d+)\\.myproj\\-([a-z]+)\\.dc\\d+\\.net")
	input := HostConfigInput {
		Hostnames: []string{
			"x-master1.myproj-prod.dc1.net",
			"x-slave1.myproj-test.dc2.net",
		},
		HostnameRegexp: r,
		TargetPatternTemplate: "%2-%1",
		HostConfig: &HostConfig{
			IdentityFile: "~/.ssh/id_rsa",
		},
	}

	// when
	results, err := NewCompiler().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "prod-master1", results[0].Host)
	assert.Equal(t, "test-slave1", results[1].Host)
}

func Test_ShouldReplaceAllGroupMatchOccurrences(t *testing.T) {
	t.Parallel()

	// given
	r, _ := regexp.Compile("x-([a-z]+\\d+)\\.myproj\\-prod\\.dc\\d+\\.net")
	input := HostConfigInput {
		Hostnames: []string{
			"x-master1.myproj-prod.dc1.net",
		},
		HostnameRegexp: r,
		TargetPatternTemplate: "%1-%1-%1",
	}

	// when
	results, err := NewCompiler().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "master1-master1-master1", results[0].Host)
}


func Test_ShouldExpandHostnameWithProvidedRange(t *testing.T) {
	t.Parallel()

	// given
	r, _ := regexp.Compile("x-([a-z]+\\d+).*")
	input := HostConfigInput {
		Hostnames: []string{
			"x-master[4..6].myproj-prod.dc1.net",
		},
		HostnameRegexp: r,
		TargetPatternTemplate: "%1",
	}

	// when
	results, err := NewCompiler().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Len(t, results, 3)
	assert.Equal(t, results[0].Host, "master4")
	assert.Equal(t, results[1].Host, "master5")
	assert.Equal(t, results[2].Host, "master6")
}