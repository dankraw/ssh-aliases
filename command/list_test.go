package command

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fixtureDir = "test-fixtures"

func TestListCommandExecute(t *testing.T) {
	t.Parallel()

	// given
	buffer := new(bytes.Buffer)
	hosts := []string{}

	// when
	err := newListCommand(buffer).execute(fixtureDir, hosts)

	// then
	assert.NoError(t, err)
	output, _ := os.ReadFile(filepath.Join(fixtureDir, "list_result"))
	assert.Equal(t, string(output), buffer.String())

}
