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
	err := NewCompileCommand(buffer).Execute(FIXTURE_DIR)

	// then
	assert.NoError(t, err)
	output, _ := ioutil.ReadFile(filepath.Join(FIXTURE_DIR, "print_compile"))
	assert.Equal(t, string(output), buffer.String())
}
