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
	files, err := scanner.ScanDirectory("../config_test/test_fixtures/valid")

	// then
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"../config_test/test_fixtures/valid/empty.hcl",
		"../config_test/test_fixtures/valid/example.hcl",
		"../config_test/test_fixtures/valid/variables.hcl",
	}, files)
}
