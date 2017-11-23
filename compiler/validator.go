package compiler

import (
	"fmt"
)

// Validator can be used to check the consistency of generated HostEntities
type Validator struct{}

// NewValidator creates an instance of Validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateResults checks if generated HostEntities have unique alias names
func (v *Validator) ValidateResults(results []HostEntity) error {
	aliases := make(map[string]struct{})
	var exists struct{}
	for _, r := range results {
		if _, contains := aliases[r.Host]; contains {
			return fmt.Errorf("generated results contain duplicate alias: `%v`", r.Host)
		}
		aliases[r.Host] = exists
	}
	return nil
}
