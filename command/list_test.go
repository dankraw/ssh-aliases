package command

import (
	"testing"

	"bytes"

	"io/ioutil"
	"path/filepath"

	"github.com/stretchr/testify/assert"
)

const fixtureDir = "test-fixtures"

func TestListCommandExecute(t *testing.T) {
	t.Parallel()

	// given
	buffer := new(bytes.Buffer)

	// when
	err := NewListCommand(buffer).Execute(fixtureDir)

	// then
	assert.NoError(t, err)
	output, _ := ioutil.ReadFile(filepath.Join(fixtureDir, "print_list"))
	assert.Equal(t, string(output), buffer.String())

}
