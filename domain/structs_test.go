package domain

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldSortHostConfig(t *testing.T) {
	t.Parallel()

	// given
	config := HostConfig{
		"b": "something",
		"c": "abc",
		"d": 0,
		"a": 123,
	}
	entries := config.ToHostConfigEntries()

	// when
	sort.Sort(ByHostConfigEntryKey(entries))

	// then
	assert.Equal(t, []HostConfigEntry{
		{"a", 123},
		{"b", "something"},
		{"c", "abc"},
		{"d", 0},
	}, entries)
}
