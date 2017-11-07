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
	assert.Equal(t, []ExpandedHostname{{
		Hostname: "x-master1.myproj-prod.dc1.net",
	}}, hostnames)
}

func TestShouldExpandHostnameWithRange(t *testing.T) {
	t.Parallel()

	// given
	hostname := "x-master[1..3].myproj-prod.dc1.net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []ExpandedHostname{{
		Hostname:     "x-master1.myproj-prod.dc1.net",
		Replacements: []string{"1"},
	}, {
		Hostname:     "x-master2.myproj-prod.dc1.net",
		Replacements: []string{"2"},
	}, {
		Hostname:     "x-master3.myproj-prod.dc1.net",
		Replacements: []string{"3"},
	}}, hostnames)
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

func TestShouldReturnErrorWhenProducedStringIsNotAValidHostname(t *testing.T) {
	t.Parallel()

	// given
	hostname := "--ddd--[1..2]..."

	// when
	_, err := NewExpander().expand(hostname)

	// then
	assert.Error(t, err)
	assert.Equal(t, "Produced string `--ddd--1...` is not a valid Hostname", err.Error())
}

func TestShouldReturnErrorWhenNoRangeWasFoundAndProducedStringIsNotAValidHostname(t *testing.T) {
	t.Parallel()

	// given
	hostname := "--ddd--[1..2..."

	// when
	_, err := NewExpander().expand(hostname)

	// then
	assert.Error(t, err)
	assert.Equal(t, "Produced string `--ddd--[1..2...` is not a valid Hostname", err.Error())
}

func TestShouldExpandHostnameWithMultipleRanges(t *testing.T) {
	t.Parallel()
	// given
	hostname := "x-master[1..3].myproj-prod.dc[1..2].net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []ExpandedHostname{{
		Hostname:     "x-master1.myproj-prod.dc1.net",
		Replacements: []string{"1", "1"},
	}, {
		Hostname:     "x-master2.myproj-prod.dc1.net",
		Replacements: []string{"2", "1"},
	}, {
		Hostname:     "x-master3.myproj-prod.dc1.net",
		Replacements: []string{"3", "1"},
	}, {
		Hostname:     "x-master1.myproj-prod.dc2.net",
		Replacements: []string{"1", "2"},
	}, {
		Hostname:     "x-master2.myproj-prod.dc2.net",
		Replacements: []string{"2", "2"},
	}, {
		Hostname:     "x-master3.myproj-prod.dc2.net",
		Replacements: []string{"3", "2"},
	}}, hostnames)
}

func TestShouldReturnErrorForSingleVariation(t *testing.T) {
	t.Parallel()

	// given
	hostname := "server-[prod].myproj.net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []ExpandedHostname{{
		Hostname:     "server-prod.myproj.net",
		Replacements: []string{"prod"},
	}}, hostnames)
}

func TestShouldAllowVariationOnBeginningOfHostname(t *testing.T) {
	t.Parallel()

	// given
	hostname := "[a|b]-server.myproj.net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []ExpandedHostname{{
		Hostname:     "a-server.myproj.net",
		Replacements: []string{"a"},
	}, {
		Hostname:     "b-server.myproj.net",
		Replacements: []string{"b"},
	}}, hostnames)
}

func TestShouldAllowVariationOnEndingOfHostname(t *testing.T) {
	t.Parallel()

	// given
	hostname := "server.myproj.[net|com]"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []ExpandedHostname{{
		Hostname:     "server.myproj.net",
		Replacements: []string{"net"},
	}, {
		Hostname:     "server.myproj.com",
		Replacements: []string{"com"},
	}}, hostnames)
}

func TestShouldExpandHostnameWithVariations(t *testing.T) {
	t.Parallel()

	// given
	hostname := "server-[prod|test|dev].myproj.net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []ExpandedHostname{{
		Hostname:     "server-prod.myproj.net",
		Replacements: []string{"prod"},
	}, {
		Hostname:     "server-test.myproj.net",
		Replacements: []string{"test"},
	}, {
		Hostname:     "server-dev.myproj.net",
		Replacements: []string{"dev"},
	}}, hostnames)
}

func TestShouldExpandHostnameWithRangesAndVariations(t *testing.T) {
	t.Parallel()

	// given
	hostname := "host[1..2].server-[prod|test].myproj[5..6].net"

	// when
	hostnames, err := NewExpander().expand(hostname)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []ExpandedHostname{{
		Hostname:     "host1.server-prod.myproj5.net",
		Replacements: []string{"1", "prod", "5"},
	}, {
		Hostname:     "host2.server-prod.myproj5.net",
		Replacements: []string{"2", "prod", "5"},
	}, {
		Hostname:     "host1.server-test.myproj5.net",
		Replacements: []string{"1", "test", "5"},
	}, {
		Hostname:     "host2.server-test.myproj5.net",
		Replacements: []string{"2", "test", "5"},
	}, {
		Hostname:     "host1.server-prod.myproj6.net",
		Replacements: []string{"1", "prod", "6"},
	}, {
		Hostname:     "host2.server-prod.myproj6.net",
		Replacements: []string{"2", "prod", "6"},
	}, {
		Hostname:     "host1.server-test.myproj6.net",
		Replacements: []string{"1", "test", "6"},
	}, {
		Hostname:     "host2.server-test.myproj6.net",
		Replacements: []string{"2", "test", "6"},
	}}, hostnames)
}
