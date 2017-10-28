package command

import (
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

const FIXTURE_DIR = "../config/test-fixtures"

func TestCompile(t *testing.T) {
	t.Parallel()

	// given

	// when
	err := NewListCommand(os.Stdout).List(FIXTURE_DIR)

	// then
	assert.NoError(t, err)
}
