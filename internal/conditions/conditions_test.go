package conditions

import (
	"testing"

	"github.com/maxim-nazarenko/tf-module-update/internal/module"
	"github.com/maxim-nazarenko/tf-module-update/internal/testhelpers"
)

func TestConditionsAll(t *testing.T) {
	testCases := []struct {
		name           string
		conditions     []Condition
		expectedResult bool
	}{
		{
			name: "all conditions resolve to true",
			conditions: []Condition{
				func(s module.Source) bool { return true },
				func(s module.Source) bool { return true },
				func(s module.Source) bool { return true },
			},
			expectedResult: true,
		},
		{
			name: "one false condition",
			conditions: []Condition{
				func(s module.Source) bool { return false },
				func(s module.Source) bool { return true },
				func(s module.Source) bool { return true },
			},
			expectedResult: false,
		},
		{
			name: "all conditions are false",
			conditions: []Condition{
				func(s module.Source) bool { return false },
				func(s module.Source) bool { return false },
				func(s module.Source) bool { return false },
			},
			expectedResult: false,
		},
	}

	var result bool
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := testhelpers.Assert(t)

			result = All(tc.conditions...)(module.Source{})
			assert.Equal(tc.expectedResult, result)
		})
	}
}

func TestConditionsAny(t *testing.T) {
	testCases := []struct {
		name           string
		conditions     []Condition
		expectedResult bool
	}{
		{
			name: "all conditions resolve to true",
			conditions: []Condition{
				func(s module.Source) bool { return true },
				func(s module.Source) bool { return true },
				func(s module.Source) bool { return true },
			},
			expectedResult: true,
		},
		{
			name: "all conditions but one are false",
			conditions: []Condition{
				func(s module.Source) bool { return false },
				func(s module.Source) bool { return true },
				func(s module.Source) bool { return false },
			},
			expectedResult: true,
		},
		{
			name: "all conditions are false",
			conditions: []Condition{
				func(s module.Source) bool { return false },
				func(s module.Source) bool { return false },
				func(s module.Source) bool { return false },
			},
			expectedResult: false,
		},
	}

	var result bool
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := testhelpers.Assert(t)

			result = Any(tc.conditions...)(module.Source{})
			assert.Equal(tc.expectedResult, result)
		})
	}
}

func TestModuleMatches(t *testing.T) {
	testCases := []struct {
		name           string
		moduleName     string
		moduleSource   module.Source
		expectedResult bool
	}{
		{
			name:           "module matches",
			moduleName:     "/example-org/aws/vpc",
			moduleSource:   module.Source{Scheme: "https", Host: "github.com", Module: "/example-org/aws/vpc", Revision: module.Revision("v1.0.0")},
			expectedResult: true,
		},
		{
			name:           "module does not match",
			moduleName:     "/aws/vpc",
			moduleSource:   module.Source{Scheme: "https", Host: "github.com", Module: "/example-org/aws/vpc_new", Revision: module.Revision("v1.0.0")},
			expectedResult: false,
		},
	}

	var result bool
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := testhelpers.Assert(t)

			result = ModuleMatches(tc.moduleName)(tc.moduleSource)
			assert.Equal(tc.expectedResult, result)
		})
	}
}

func TestSubmoduleMatches(t *testing.T) {
	testCases := []struct {
		name           string
		submoduleName  string
		moduleSource   module.Source
		expectedResult bool
	}{
		{
			name:           "submodule matches",
			submoduleName:  "vpc_enhanced",
			moduleSource:   module.Source{Scheme: "https", Host: "github.com", Submodule: "vpc_enhanced", Module: "/example-org/aws/vpc", Revision: module.Revision("v1.0.0")},
			expectedResult: true,
		},
		{
			name:          "submodule does not match",
			submoduleName: "vpc_old",
			moduleSource:  module.Source{Scheme: "https", Host: "github.com", Submodule: "vpc_enhanced", Module: "/example-org/aws/vpc", Revision: module.Revision("v1.0.0")},

			expectedResult: false,
		},
	}

	var result bool
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := testhelpers.Assert(t)

			result = SubmoduleMatches(tc.submoduleName)(tc.moduleSource)
			assert.Equal(tc.expectedResult, result)
		})
	}
}

func TestRevisionMatches(t *testing.T) {
	testCases := []struct {
		name           string
		revision       module.Revision
		moduleSource   module.Source
		expectedResult bool
	}{
		{
			name:           "revision matches",
			revision:       module.Revision("v1.2.0"),
			moduleSource:   module.Source{Scheme: "https", Host: "github.com", Module: "/example-org/aws/vpc", Revision: module.Revision("v1.2.0")},
			expectedResult: true,
		},
		{
			name:         "submodule does not match",
			revision:     module.Revision("v2.0.1"),
			moduleSource: module.Source{Scheme: "https", Host: "github.com", Submodule: "vpc_enhanced", Module: "/example-org/aws/vpc", Revision: module.Revision("v1.0.0")},

			expectedResult: false,
		},
	}

	var result bool
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := testhelpers.Assert(t)

			result = RevisionMatches(tc.revision)(tc.moduleSource)
			assert.Equal(tc.expectedResult, result)
		})
	}
}
