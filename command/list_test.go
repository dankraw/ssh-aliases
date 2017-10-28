package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const FIXTURE_DIR = "../config/test-fixtures"

func TestCompile(t *testing.T) {
	t.Parallel()

	// when
	err := NewListCommand().List(FIXTURE_DIR)

	// then
	assert.NoError(t, err)
}
