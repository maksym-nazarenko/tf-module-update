package conditions

import "github.com/maxim-nazarenko/tf-module-update/internal/module"

// Condition is binary function to make decision
type Condition func(module.Source) bool
