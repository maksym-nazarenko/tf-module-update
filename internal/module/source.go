package module

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// String constructs string representation of the module source
//
// Basically, the following should be always true:
//
// var sourceURL string = "..."
// source, _ := ParseSource(sourceURL)
// source.String() == sourceURL
//
func (s Source) String() string {
	revision := ""
	if s.Revision != "" {
		revision = "?ref=" + string(s.Revision)
	}

	scheme := ""
	if s.Scheme != "" {
		scheme = s.Scheme + "://"
	}

	return s.SpecialPrefix + scheme + s.Host + s.Module + s.Submodule + revision
}

// Merge combines two sources and returns new struct
// This function overrides fields in calling struct with fields from other object
// but only if the incoming field is not empty
func (s Source) Merge(o Source) Source {
	merged := s

	if o.Scheme != "" {
		merged.Scheme = o.Scheme
	}

	if o.Host != "" {
		merged.Host = o.Host
	}

	if o.Module != "" {
		merged.Module = o.Module
	}

	if o.Submodule != "" {
		merged.Submodule = o.Submodule
	}

	if o.Revision != "" {
		merged.Revision = o.Revision
	}

	return merged
}

// ParseSource builds struct from string representation of module source
func ParseSource(source string) (Source, error) {
	result := Source{}

	if len(source) == 0 {
		return result, nil
	}

	// special case for well-known hostname
	if strings.HasPrefix(source, "github.com") {
		source = "https://" + source
	}

	var specialPrefix string = ""
	if strings.Contains(source, "::") {
		specialPrefix = source[:strings.LastIndex(source, "::")+2]
	}

	// weird source which contains only special prefix
	if len(specialPrefix) == len(source) {
		return result, fmt.Errorf("module source consists only of special prefix")
	}

	// only "git::" special prefix is supported at the moment
	if len(specialPrefix) > 0 && specialPrefix != "git::" {
		return result, &InvalidSourceFormatError{fmt.Sprintf("only 'git::' special prefix is supported but got '%s', skipping", specialPrefix)}
	}

	if len(specialPrefix) > 0 {
		source = source[len(specialPrefix):]
	}

	parsedSource, err := url.Parse(source)
	if err != nil {
		return result, errors.New("cannot parse source string: " + err.Error())
	}

	modulePath := parsedSource.Path
	submodule := ""
	if strings.Contains(modulePath, "//") {
		submodule = modulePath[strings.Index(modulePath, "//"):]
		modulePath = modulePath[:strings.Index(modulePath, "//")]
	}

	query := parsedSource.Query()

	// there are more than 1 query parameter
	// but we are not aware of them, so play defense
	// to not break anything
	if len(query) > 1 {
		return result, &InvalidSourceFormatError{"more than 1 query parameter found"}
	}

	ref := query.Get("ref")
	if len(query) == 1 && ref == "" {
		return result, &InvalidSourceFormatError{"query param is provided but it is not 'ref'"}
	}

	return Source{
		Scheme:        parsedSource.Scheme,
		Host:          parsedSource.Host,
		SpecialPrefix: specialPrefix,
		Module:        modulePath,
		Submodule:     submodule,
		Revision:      Revision(ref),
	}, nil
}
