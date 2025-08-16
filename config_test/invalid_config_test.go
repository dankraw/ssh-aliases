package config

import (
	"path/filepath"
	"testing"

	"github.com/dankraw/ssh-aliases/config"
	"github.com/stretchr/testify/assert"
)

var reader = config.NewReader()

var testsParentDir = filepath.Join("test_fixtures", "invalid")

var tests = []struct {
	dir              string
	expectedErrorMsg string
}{
	{"duplicate_alias", "duplicate host `service-a`"},
	{"config_not_found", "error in `test_fixtures/invalid/config_not_found/host_only.hcl`: " +
		"error in `wally-host` host definition: no config `wally` found"},
	{"invalid_hcl", "failed parsing `test_fixtures/invalid/invalid_hcl/invalid.hcl`: " +
		"At 7:2: object expected closing RBRACE got: EOF"},
	{"variable_redeclaration", "error in `test_fixtures/invalid/variable_redeclaration/example.hcl`: " +
		"variable redeclaration: `abc.def`"},
	{"invalid_import_value", "error in `test_fixtures/invalid/invalid_import_value/example.hcl`: " +
		"invalid `def_conf` config definition: config import statement has invalid value: `1`"},
	{"alias_and_hostname_not_specified", "error in `test_fixtures/invalid/alias_and_hostname_not_specified/example.hcl`: " +
		"invalid `wat` host definition: alias and hostname are both empty or undefined"},
	{"no_hostname_nor_config", "error in `test_fixtures/invalid/no_hostname_nor_config/example.hcl`: " +
		"no config nor hostname specified for host `wat`"},
	{"non_existing_variable/in_alias", "error in `test_fixtures/invalid/non_existing_variable/in_alias/example.hcl`: " +
		"error in alias of `service-a` host definition: variable `b.c3.d4` not defined"},
	{"non_existing_variable/in_hostname", "error in `test_fixtures/invalid/non_existing_variable/in_hostname/example.hcl`: " +
		"error in hostname of `service-a` host definition: variable `b.c3ę.d4.E_F-d` not defined"},
	{"non_existing_variable/in_config", "error in `test_fixtures/invalid/non_existing_variable/in_config/example.hcl`: " +
		"error in `service-a` host definition: could not compile config property `user`: variable `b.c3.d4` not defined"},
	{"non_existing_variable/in_external_config", "error in `test_fixtures/invalid/non_existing_variable/in_external_config/example.hcl`: " +
		"invalid `ext` config definition: could not compile config property `user`: variable `b.c3.d4` not defined"},
}

func TestShouldThrowErrorOnDuplicateAlias(t *testing.T) {
	t.Parallel()

	for _, test := range tests {
		t.Run(test.dir, func(t *testing.T) {
			t.Parallel()
			// when
			_, err := reader.ReadConfigs(filepath.Join(testsParentDir, test.dir))

			// then
			assert.Error(t, err)
			assert.Equal(t, test.expectedErrorMsg, err.Error())
		})
	}
}

func TestShouldThrowErrorOnCircularImports(t *testing.T) {
	t.Parallel()

	// when
	_, err := reader.ReadConfigs("./test_fixtures/invalid/circular_imports")

	// then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular import in configs")
}
