package strategies

import (
	"github.com/maxim-nazarenko/tf-module-update/internal/conditions"
	"github.com/maxim-nazarenko/tf-module-update/internal/module"
)

// MutatorFunc is type to change module source
type MutatorFunc func(module.Source) module.Source

// Strict strategy checks that all conditions are satisfied
//
// If strategy decision is positive then sourceMutator function is applied to module source
type Strict struct {
	conditions    []conditions.Condition
	sourceMutator MutatorFunc
}

// WithCondition adds condition to the chain of conditions
func (u *Strict) WithCondition(cond conditions.Condition) *Strict {
	u.conditions = append(u.conditions, cond)

	return u
}

// WithConditions is a convenient way to add multiple conditions at once
func (u *Strict) WithConditions(conds ...conditions.Condition) *Strict {
	u.conditions = append(u.conditions, conds...)

	return u
}

// Apply creates updated clone of module source using sourceMutator function
func (u *Strict) Apply(source module.Source) module.Source {
	return u.sourceMutator(source)
}

// Decide checks all conditions to make decision if the module source should be updated
// Returns false if no conditions were applied during checking
func (u *Strict) Decide(source module.Source) bool {
	if len(u.conditions) == 0 {
		return false
	}

	for _, v := range u.conditions {
		if !v(source) {
			return false
		}
	}
	return true
}

func NewStrictUpdater(sourceMutator MutatorFunc) *Strict {
	return &Strict{sourceMutator: sourceMutator}
}
