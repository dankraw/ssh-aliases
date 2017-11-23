package examples

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/dankraw/ssh-aliases/command"
	"github.com/stretchr/testify/assert"
)

const dir = "readme"

func TestCompileCommandExecute(t *testing.T) {
	t.Parallel()

	// given
	buffer := new(bytes.Buffer)

	// when
	cli, err := command.NewCLI("test-version", buffer)

	// then
	assert.NoError(t, err)

	// and
	err = cli.ApplyArgs([]string{"ssh-aliases", "--scan", dir, "compile"})

	// then
	assert.NoError(t, err)
	output, _ := ioutil.ReadFile(filepath.Join(dir, "compile_result"))
	assert.Equal(t, string(output), buffer.String())
}

func TestListCommandExecute(t *testing.T) {
	t.Parallel()

	// given
	buffer := new(bytes.Buffer)

	// when
	cli, err := command.NewCLI("test-version", buffer)

	// then
	assert.NoError(t, err)

	// and
	err = cli.ApplyArgs([]string{"ssh-aliases", "--scan", dir, "list"})

	// then
	assert.NoError(t, err)
	output, _ := ioutil.ReadFile(filepath.Join(dir, "list_result"))
	assert.Equal(t, string(output), buffer.String())
}
