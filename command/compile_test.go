package command

import (
	"bytes"
	"testing"

	"io/ioutil"
	"path/filepath"

	"github.com/stretchr/testify/assert"
)

func TestCompileCommandExecute(t *testing.T) {
	t.Parallel()

	// given
	buffer := new(bytes.Buffer)

	// when
	err := NewCompileCommand(buffer).Execute(fixtureDir)

	// then
	assert.NoError(t, err)
	output, _ := ioutil.ReadFile(filepath.Join(fixtureDir, "print_compile"))
	assert.Equal(t, string(output), buffer.String())
}
