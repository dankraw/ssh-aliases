package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldScanDir(t *testing.T) {
	t.Parallel()

	// given
	scanner := NewScanner()

	// when
	hcls, err := scanner.ScanDirectory(fixtureDir)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"./test-fixtures/empty.hcl",
		"./test-fixtures/example.hcl",
	}, hcls)
}
