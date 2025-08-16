package command

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompileCommandExecute(t *testing.T) {
	t.Parallel()

	// given
	buffer := new(bytes.Buffer)
	hosts := []string{}

	// when
	err := newCompileCommand(buffer).execute(fixtureDir, hosts)

	// then
	assert.NoError(t, err)
	output, _ := os.ReadFile(filepath.Join(fixtureDir, "compile_result"))
	assert.Equal(t, string(output), buffer.String())
}
