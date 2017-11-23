package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldScanDir(t *testing.T) {
	t.Parallel()

	// given
	scanner := NewScanner()

	// when
	files, err := scanner.ScanDirectory(fixtureDir)

	// then
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"test-fixtures/empty.hcl",
		"test-fixtures/example.hcl",
	}, files)
}
