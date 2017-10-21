package compiler

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_ShouldExpandHostname(t *testing.T) {
	t.Parallel()
	// given
	hostname := "x-master[1..3].myproj-prod.dc1.net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"x-master1.myproj-prod.dc1.net",
		"x-master2.myproj-prod.dc1.net",
		"x-master3.myproj-prod.dc1.net",
	}, hostnames)
}

func Test_ShouldExpandHostnameWithMultipleRanges(t *testing.T) {
	t.Parallel()
	// given
	hostname := "x-master[1..3].myproj-prod.dc[1..2].net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"x-master1.myproj-prod.dc1.net",
		"x-master1.myproj-prod.dc2.net",
		"x-master2.myproj-prod.dc1.net",
		"x-master2.myproj-prod.dc2.net",
		"x-master3.myproj-prod.dc1.net",
		"x-master3.myproj-prod.dc2.net",
	}, hostnames)
}
