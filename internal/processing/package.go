package processing

import (
	"io/fs"
	"strings"
)

// DefaultExclusionNamesList contains well-known file and folder names to be excluded from processing
var DefaultExclusionNamesList []string = []string{
	".terraform",
	".git",
}

// DefaultExclusionFunc excludes all hidden, e.g. starting with dot ".", folders and files
var DefaultExclusionFunc ExcludeFileFunc = func(info fs.FileInfo) bool {
	return !strings.HasPrefix(info.Name(), ".")
}
