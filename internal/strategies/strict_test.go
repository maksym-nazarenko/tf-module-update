package strategies

import (
	"testing"

	"github.com/maxim-nazarenko/tf-module-update/internal/conditions"
	"github.com/maxim-nazarenko/tf-module-update/internal/module"
	"github.com/maxim-nazarenko/tf-module-update/internal/testhelpers"
)

func TestStrictDecide(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult bool
		conditions     []conditions.Condition
		moduleSource   module.Source
	}{
		{
			name:           "no conditions added to strategy should be treated as negative decision",
			expectedResult: false,
			conditions:     []conditions.Condition{},
			moduleSource:   module.Source{Scheme: "https", Host: "example.com", Module: "/aws/vpc", Revision: module.Revision("v1.2.3")},
		},
		{
			name:           "one failed condition fails the whole chain",
			expectedResult: false,
			conditions: []conditions.Condition{
				func(s module.Source) bool { return true },
				func(s module.Source) bool { return true },
				func(s module.Source) bool { return false },
				func(s module.Source) bool { return true },
			},
			moduleSource: module.Source{Scheme: "https", Host: "example.com", Module: "/aws/vpc", Revision: module.Revision("v1.2.3")},
		},
		{
			name:           "all conditions true",
			expectedResult: true,
			conditions: []conditions.Condition{
				func(s module.Source) bool { return true },
				func(s module.Source) bool { return true },
				func(s module.Source) bool { return true },
			},
			moduleSource: module.Source{Scheme: "https", Host: "example.com", Module: "/aws/vpc", Revision: module.Revision("v1.2.3")},
		},
	}
	assert := testhelpers.Assert(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strategy := Strict{}
			strategy.WithConditions(tc.conditions...)
			assert.Equal(tc.expectedResult, strategy.Decide(tc.moduleSource))
		})
	}
}
