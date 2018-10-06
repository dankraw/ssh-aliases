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
	assert.Equal(t, "error in `test_fixtures/invalid/config_not_found/host_only.hcl`: error in `wally-host` host definition: no config `wally` found", err.Error())
}

func TestShouldThrowErrorOnInvalidHcl(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/invalid_hcl")

	// then
	assert.Error(t, err)
	assert.Equal(t, "failed parsing `test_fixtures/invalid/invalid_hcl/invalid.hcl`: At 7:2: object expected closing RBRACE got: EOF", err.Error())
}

func TestShouldThrowErrorOnValueRedeclaration(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/variable_redeclaration")

	// then
	assert.Error(t, err)
	assert.Equal(t, "error in `test_fixtures/invalid/variable_redeclaration/example.hcl`: variable redeclaration: `abc.def`", err.Error())
}

func TestShouldThrowErrorOnCircularImports(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/circular_imports")

	// then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular import in configs")
}

func TestShouldThrowErrorOnInvalidImportValue(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/invalid_import_value")

	// then
	assert.Error(t, err)
	assert.Equal(t, "error in `test_fixtures/invalid/invalid_import_value/example.hcl`: invalid `def_conf` config definition: config import statement has invalid value: `1`", err.Error())
}

func TestShouldThrowErrorOnNoAliasSpecified(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/alias_and_hostname_not_specified")

	// then
	assert.Error(t, err)
	assert.Equal(t, "error in `test_fixtures/invalid/alias_and_hostname_not_specified/example.hcl`: invalid `wat` host definition: alias and hostname are both empty or undefined", err.Error())
}

func TestShouldThrowErrorOnNoHostnameNorConfigSpecified(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/no_hostname_nor_config")

	// then
	assert.Error(t, err)
	assert.Equal(t, "error in `test_fixtures/invalid/no_hostname_nor_config/example.hcl`: no config nor hostname specified for host `wat`", err.Error())
}

func TestShouldThrowErrorOnNonExistingVariableInAlias(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/non_existing_variable/in_alias")

	// then
	assert.Error(t, err)
	assert.Equal(t, "error in `test_fixtures/invalid/non_existing_variable/in_alias/example.hcl`: error in alias of `service-a` host definition: variable `b.c3.d4` not defined", err.Error())
}

func TestShouldThrowErrorOnNonExistingVariableInHostname(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/non_existing_variable/in_hostname")

	// then
	assert.Error(t, err)
	assert.Equal(t, "error in `test_fixtures/invalid/non_existing_variable/in_hostname/example.hcl`: error in hostname of `service-a` host definition: variable `b.c3ę.d4.E_F-d` not defined", err.Error())
}

func TestShouldThrowErrorOnNonExistingVariableInConfig(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/non_existing_variable/in_config")

	// then
	assert.Error(t, err)
	assert.Equal(t, "error in `test_fixtures/invalid/non_existing_variable/in_config/example.hcl`: error in `service-a` host definition: could not compile config property `user`: variable `b.c3.d4` not defined", err.Error())
}

func TestShouldThrowErrorOnNonExistingVariableInExternalConfig(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/non_existing_variable/in_external_config")

	// then
	assert.Error(t, err)
	assert.Equal(t, "error in `test_fixtures/invalid/non_existing_variable/in_external_config/example.hcl`: invalid `ext` config definition: could not compile config property `user`: variable `b.c3.d4` not defined", err.Error())
}
