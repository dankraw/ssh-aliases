package config

import (
	"testing"

	"github.com/dankraw/ssh-aliases/config"
	"github.com/stretchr/testify/assert"
)

var reader = config.NewReader()

func TestShouldThrowErrorOnDuplicateAlias(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/duplicate_alias")

	// then
	assert.Error(t, err)
	assert.Equal(t, "duplicate host `service-a`", err.Error())
}

func TestShouldThrowErrorOnNotFoundConfig(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/config_not_found")

	// then
	assert.Error(t, err)
	assert.Equal(t, "no config `wally` found (used by host `wally-host`)", err.Error())
}

func TestShouldThrowErrorOnInvalidHcl(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/invalid_hcl")

	// then
	assert.Error(t, err)
	assert.Equal(t, "failed parsing test_fixtures/invalid/invalid_hcl/invalid.hcl: At 7:2: object expected closing RBRACE got: EOF", err.Error())
}

func TestShouldThrowErrorOnValueRedeclaration(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/value_redeclaration")

	// then
	assert.Error(t, err)
	assert.Equal(t, "value redeclaration: abc.def", err.Error())
}
