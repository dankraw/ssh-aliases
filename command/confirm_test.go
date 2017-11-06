package command

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	t.Parallel()

	// when
	exists, err := NewConfirm(os.Stdin).fileExists(filepath.Join(fixtureDir, "print_list"))

	// then
	assert.True(t, exists)
	assert.NoError(t, err)
}

func TestFileNotExists(t *testing.T) {
	t.Parallel()

	// when
	exists, err := NewConfirm(os.Stdin).fileExists(filepath.Join(fixtureDir, "not_exists"))

	// then
	assert.False(t, exists)
	assert.NoError(t, err)
}

func TestDirExists(t *testing.T) {
	t.Parallel()

	// when
	exists, err := NewConfirm(os.Stdin).fileExists(fixtureDir)

	// then
	assert.True(t, exists)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("Path %s is a directory", fixtureDir), err.Error())
}

func TestInvalidPath(t *testing.T) {
	t.Parallel()

	// when
	exists, err := NewConfirm(os.Stdin).fileExists("")

	// then
	assert.False(t, exists)
	assert.Error(t, err)
	assert.Equal(t, "Provided path is empty", err.Error())
}

type TestReader struct {
	response string
}

func NewTestReader(response string) *TestReader {
	return &TestReader{
		response: response,
	}
}

func (r *TestReader) Read(p []byte) (n int, err error) {
	copy(p[:], r.response)
	return len(r.response), nil
}

func TestConfirm(t *testing.T) {
	t.Parallel()

	// when
	reader := NewTestReader("Y\n")
	confirmed, err := NewConfirm(reader).RequireConfirmationIfFileExists(filepath.Join(fixtureDir, "print_list"))

	// then
	assert.True(t, confirmed)
	assert.NoError(t, err)
}

func TestGiveUp(t *testing.T) {
	t.Parallel()

	// when
	reader := NewTestReader("nope\n")
	confirmed, err := NewConfirm(reader).RequireConfirmationIfFileExists(filepath.Join(fixtureDir, "print_list"))

	// then
	assert.False(t, confirmed)
	assert.NoError(t, err)
}
