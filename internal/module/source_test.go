package module

import (
	"testing"

	"github.com/maxim-nazarenko/tf-module-update/internal/testhelpers"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		sourceString   string
		expectedStruct Source
	}{
		{
			name:           "https proto for github.com",
			expectedError:  nil,
			sourceString:   "https://github.com/example-org/aws/vpc.git?ref=0.0.1",
			expectedStruct: Source{Scheme: "https", Host: "github.com", Module: "/example-org/aws/vpc.git", Revision: Revision("0.0.1")},
		},
		{
			name:           "git:: prefix handling",
			expectedError:  nil,
			sourceString:   "git::https://example.com/example-org/aws/vpc.git?ref=0.0.1",
			expectedStruct: Source{Scheme: "https", SpecialPrefix: "git::", Host: "example.com", Module: "/example-org/aws/vpc.git", Revision: Revision("0.0.1")},
		},
		{
			name:           "github.com without scheme defaults to https",
			expectedError:  nil,
			sourceString:   "github.com/example-org/aws/vpc.git?ref=0.0.2",
			expectedStruct: Source{Scheme: "https", Host: "github.com", Module: "/example-org/aws/vpc.git", Revision: Revision("0.0.2")},
		},
		{
			name:           "submodule",
			expectedError:  nil,
			sourceString:   "https://github.com/example-org/aws/vpc.git//src/multizone?ref=0.0.2",
			expectedStruct: Source{Scheme: "https", Host: "github.com", Module: "/example-org/aws/vpc.git", Submodule: "//src/multizone", Revision: Revision("0.0.2")},
		},
		{
			name:           "empty string returns empty struct",
			expectedError:  nil,
			sourceString:   "",
			expectedStruct: Source{},
		},
		// negative scenarios
		{
			name:           "error if no revision in source URL",
			expectedError:  &InvalidSourceFormatError{},
			sourceString:   "github.com/example-org/aws/vpc.git",
			expectedStruct: Source{},
		},
		{
			name:           "unsupported prefix",
			expectedError:  &InvalidSourceFormatError{},
			sourceString:   "hg::example.com/example-org/aws/vpc.git?ref=0.0.1",
			expectedStruct: Source{},
		},
		{
			name:           "more than one query parameter is unsupported",
			expectedError:  &InvalidSourceFormatError{},
			sourceString:   "git::example.com/example-org/aws/vpc.git?ref=0.0.1&new=true",
			expectedStruct: Source{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := testhelpers.Assert(t)

			result, err := ParseSource(tc.sourceString)

			if tc.expectedError != nil {
				assert.SameType(tc.expectedError, err)
				return
			}
			assert.NoError(err)
			assert.Equal(tc.expectedStruct, result)
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult string
		sourceStruct   Source
	}{
		{
			name:           "https github.com",
			expectedResult: "https://github.com/example-org/aws/vpc.git?ref=0.0.1",
			sourceStruct:   Source{Scheme: "https", Host: "github.com", Module: "/example-org/aws/vpc.git", Revision: Revision("0.0.1")},
		},
		{
			name:           "git:: special prefix example.com",
			expectedResult: "git::example.com/example-org/aws/vpc.git?ref=0.0.1",
			sourceStruct:   Source{Scheme: "", SpecialPrefix: "git::", Host: "example.com", Module: "/example-org/aws/vpc.git", Revision: Revision("0.0.1")},
		},
		{
			name:           "submodule is used",
			expectedResult: "https://example.com/example-org/aws/vpc.git//src/multizone?ref=0.0.1",
			sourceStruct:   Source{Scheme: "https", Host: "example.com", Module: "/example-org/aws/vpc.git", Submodule: "//src/multizone", Revision: Revision("0.0.1")},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := testhelpers.Assert(t)

			result := tc.sourceStruct.String()

			assert.Equal(tc.expectedResult, result)
		})
	}
}

func TestMerge(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult Source
		original       Source
		other          Source
	}{
		{
			name: "empty other",
			expectedResult: Source{
				Scheme:    "https",
				Host:      "github.com",
				Module:    "/modules/azure.git",
				Submodule: "/src/db",
				Revision:  Revision("v1.2.1"),
			},
			original: Source{
				Scheme:    "https",
				Host:      "github.com",
				Module:    "/modules/azure.git",
				Submodule: "/src/db",
				Revision:  Revision("v1.2.1"),
			},
			other: Source{},
		},
		{
			name: "empty original",
			expectedResult: Source{
				Scheme:    "https",
				Host:      "github.com",
				Module:    "/modules/azure.git",
				Submodule: "/src/db",
				Revision:  Revision("v1.2.1"),
			},
			other: Source{
				Scheme:    "https",
				Host:      "github.com",
				Module:    "/modules/azure.git",
				Submodule: "/src/db",
				Revision:  Revision("v1.2.1"),
			},
			original: Source{},
		},
		{
			name:           "both empty",
			expectedResult: Source{},
			original:       Source{},
			other:          Source{},
		},
		{
			name: "one field override",
			expectedResult: Source{
				Scheme:    "https",
				Host:      "github.com",
				Module:    "/modules/azure.git",
				Submodule: "/src/db",
				Revision:  Revision("v2.2.1"),
			},
			original: Source{
				Scheme:    "https",
				Host:      "github.com",
				Module:    "/modules/azure.git",
				Submodule: "/src/db",
				Revision:  Revision("v1.2.1"),
			},
			other: Source{Revision: Revision("v2.2.1")},
		},
		{
			name: "all fields override",
			expectedResult: Source{
				Scheme:    "ssh",
				Host:      "example.com",
				Module:    "/terraform-modules/azure.git",
				Submodule: "/db",
				Revision:  Revision("v2.2.1"),
			},
			original: Source{
				Scheme:    "https",
				Host:      "github.com",
				Module:    "/modules/azure.git",
				Submodule: "/src/db",
				Revision:  Revision("v1.2.1"),
			},
			other: Source{
				Scheme:    "ssh",
				Host:      "example.com",
				Module:    "/terraform-modules/azure.git",
				Submodule: "/db",
				Revision:  Revision("v2.2.1"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.original.Merge(tc.other)
			assert := testhelpers.Assert(t)
			assert.Equal(tc.expectedResult, result)
		})
	}
}
