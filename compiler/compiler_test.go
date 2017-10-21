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
	results, err := New().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, results[0].Host, "prod-master1")
	assert.Equal(t, results[1].Host, "test-slave1")
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
		HostConfig: &HostConfig{
			IdentityFile: "~/.ssh/id_rsa",
		},
	}

	// when
	results, err := New().Compile(input)

	// then
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, results[0].Host, "master1-master1-master1")
}