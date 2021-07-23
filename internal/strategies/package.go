package strategies

import "github.com/maxim-nazarenko/tf-module-update/internal/module"

// Strategy is a type to make decision and mutate module source string
type Strategy interface {
	Decide(module.Source) bool
	Apply(module.Source) module.Source
}
