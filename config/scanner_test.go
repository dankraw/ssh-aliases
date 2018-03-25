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
	files, err := scanner.ScanDirectory("../config_test/test_fixtures/valid/basic_with_variables")

	// then
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"../config_test/test_fixtures/valid/basic_with_variables/empty.hcl",
		"../config_test/test_fixtures/valid/basic_with_variables/example.hcl",
		"../config_test/test_fixtures/valid/basic_with_variables/variables.hcl",
	}, files)
}
