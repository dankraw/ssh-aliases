package compiler

import (
	"errors"
	"fmt"

	. "github.com/dankraw/ssh-aliases/domain"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateResults(results []HostConfigResult) error {
	aliases := make(map[string]struct{})
	var exists struct{}
	for _, r := range results {
		if _, contains := aliases[r.Host]; contains {
			return errors.New(fmt.Sprintf("Generated results contain duplicate alias: %v", r.Host))
		}
		aliases[r.Host] = exists
	}
	return nil
}
