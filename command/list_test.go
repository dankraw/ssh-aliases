package command

import (
	"testing"

	"bytes"

	"io/ioutil"
	"path/filepath"

	"github.com/stretchr/testify/assert"
)

const FIXTURE_DIR = "test-fixtures"

func TestListCommandExecute(t *testing.T) {
	t.Parallel()

	// given
	buffer := new(bytes.Buffer)

	// when
	err := NewListCommand(buffer).Execute(FIXTURE_DIR)

	// then
	assert.NoError(t, err)
	output, _ := ioutil.ReadFile(filepath.Join(FIXTURE_DIR, "print_list"))
	assert.Equal(t, string(output), buffer.String())

}
