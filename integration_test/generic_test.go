package examples

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/dankraw/ssh-aliases/command"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	dir            string
	additionalArgs []string
}{
	{"hostname_only", []string{}},
	{"readme", []string{}},
	{"readme_regexp", []string{
		"--hosts-file", filepath.Join("readme_regexp", "hosts.txt"),
	}},
	{"regexp_hosts", []string{
		"--hosts-file", filepath.Join("regexp_hosts", "hosts.txt"),
	}},
}

func TestCompileCommandExecute(t *testing.T) {
	t.Parallel()

	for _, test := range tests {
		t.Run(test.dir, func(t *testing.T) {
			// given
			buffer := new(bytes.Buffer)

			// when
			cli, err := command.NewCLI("test-version", buffer)

			// then
			assert.NoError(t, err)

			// and
			args := append([]string{"ssh-aliases", "--scan", test.dir, "compile"}, test.additionalArgs...)
			err = cli.ApplyArgs(args)

			// then
			assert.NoError(t, err)
			output, _ := ioutil.ReadFile(filepath.Join(test.dir, "compile_result"))
			assert.Equal(t, string(output), buffer.String())
		})
	}
}

func TestListCommandExecute(t *testing.T) {
	t.Parallel()

	for _, test := range tests {
		t.Run(test.dir, func(t *testing.T) {
			// given
			buffer := new(bytes.Buffer)

			// when
			cli, err := command.NewCLI("test-version", buffer)

			// then
			assert.NoError(t, err)

			// and
			args := append([]string{"ssh-aliases", "--scan", test.dir, "list"}, test.additionalArgs...)
			err = cli.ApplyArgs(args)

			// then
			assert.NoError(t, err)
			output, _ := ioutil.ReadFile(filepath.Join(test.dir, "list_result"))
			assert.Equal(t, string(output), buffer.String())
		})
	}
}
