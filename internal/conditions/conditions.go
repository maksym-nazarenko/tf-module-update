package conditions

import "github.com/maxim-nazarenko/tf-module-update/internal/module"

// All builds a composite condition that requires all conditions to return true
func All(conditions ...Condition) Condition {
	return func(s module.Source) bool {
		for _, v := range conditions {
			if !v(s) {
				return false
			}
		}

		return true
	}
}

// Any builds a composite condition that requires at least one condition to return true
func Any(conditions ...Condition) Condition {
	return func(s module.Source) bool {
		for _, v := range conditions {
			if v(s) {
				return true
			}
		}

		return false
	}
}

// ModuleMatches builds condition that returns true if module matches given module name
func ModuleMatches(moduleName string) Condition {
	return func(s module.Source) bool {
		return s.Module == moduleName
	}
}

// SubmoduleMatches builds condition that returns true if submodule matches given submodule name
func SubmoduleMatches(submoduleName string) Condition {
	return func(s module.Source) bool {
		return s.Submodule == submoduleName
	}
}

// RevisionMatches builds condition that returns true if revision matches given revision
func RevisionMatches(r module.Revision) Condition {
	return func(s module.Source) bool {
		return string(s.Revision) == string(r)
	}
}

// HostMatches builds condition that returns true if host matches given host
func HostMatches(host string) Condition {
	return func(s module.Source) bool {
		return s.Host == host
	}
}

// SchemeMatches builds condition that returns true if scheme matches given scheme
func SchemeMatches(scheme string) Condition {
	return func(s module.Source) bool {
		return s.Scheme == scheme
	}
}

// False builds condition that always returns false
func False() Condition {
	return func(s module.Source) bool {
		return false
	}
}
