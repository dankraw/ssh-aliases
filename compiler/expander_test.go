package compiler

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestShouldNotExpandForNoOperators(t *testing.T) {
	t.Parallel()
	// given
	hostname := "x-master1.myproj-prod.dc1.net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"x-master1.myproj-prod.dc1.net",
	}, hostnames)
}

func TestShouldExpandHostnameWithRange(t *testing.T) {
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

func TestShouldReturnErrorOnInvalidRange(t *testing.T) {
	t.Parallel()
	// given
	hostname := "x-master[120..13].myproj-prod.dc1.net"

	// when
	_, err := NewExpander().expand(hostname)

	// then
	assert.Error(t, err)
	assert.Equal(t, "Invalid range: 120 is not smaller than 13", err.Error())
}

func TestShouldExpandHostnameWithMultipleRanges(t *testing.T) {
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

func TestShouldExpandHostnameWithVariations(t *testing.T) {
	t.Parallel()
	// given
	hostname := "server-[prod|test|dev].myproj.net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"server-prod.myproj.net",
		"server-test.myproj.net",
		"server-dev.myproj.net",
	}, hostnames)
}

func TestShouldExpandHostnameWithRangesAndVariations(t *testing.T) {
	t.Parallel()
	// given
	hostname := "host[1..2].server-[prod|test].myproj.net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"host1.server-prod.myproj.net",
		"host1.server-test.myproj.net",
		"host2.server-prod.myproj.net",
		"host2.server-test.myproj.net",
	}, hostnames)
}

