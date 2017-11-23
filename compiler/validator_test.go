package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	// given
	results := []HostEntity{{
		Host: "is_unique",
	}, {
		Host: "is_unique",
	}}

	// when
	err := NewValidator().ValidateResults(results)

	// then
	assert.Error(t, err)
	assert.Equal(t, "generated results contain duplicate alias: `is_unique`", err.Error())
}
