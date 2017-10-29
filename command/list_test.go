package command

import (
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
)

const FIXTURE_DIR = "../config/test-fixtures"

func TestCompile(t *testing.T) {
	t.Parallel()

	// given
	buffer := new(bytes.Buffer)

	// when
	err := NewListCommand(buffer).List(FIXTURE_DIR)

	// then
	assert.NoError(t, err)
	assert.Equal(t, `../config/test-fixtures/empty.hcl (definitions=0):
../config/test-fixtures/example.hcl (definitions=2):
 service-a (compiled=5):
  a1: service-a1.example.com
  a2: service-a2.example.com
  a3: service-a3.example.com
  a4: service-a4.example.com
  a5: service-a5.example.com
 service-b (compiled=2):
  b1: service-b1.example.com
  b2: service-b2.example.com
`, buffer.String())

}
