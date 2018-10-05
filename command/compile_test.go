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
	hosts := []string{}

	// when
	err := newCompileCommand(buffer).execute(fixtureDir, hosts)

	// then
	assert.NoError(t, err)
	output, _ := ioutil.ReadFile(filepath.Join(fixtureDir, "compile_result"))
	assert.Equal(t, string(output), buffer.String())
}
