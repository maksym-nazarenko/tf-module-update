package module

// Source describes module source with possible submodule and revision
type Source struct {
	Scheme        string
	Host          string
	SpecialPrefix string
	Module        string // module name, including organization name for github
	Submodule     string
	Revision      Revision
}

// Revision represents revision of a module
type Revision string
