package testhelpers

import (
	"reflect"
	"testing"
)

type Asserter interface {
	Equal(expected, actual interface{})
}

type assertion struct {
	t *testing.T
}

var _ Asserter = (*assertion)(nil)

func (a *assertion) Equal(expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		a.t.Fatalf(`
		values are not equal:
		expected: %#v (%T)
		got: %#v (%T)
		`,
			expected, expected,
			actual, actual,
		)
	}
}

func (a *assertion) SameType(expected, actual interface{}) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		a.t.Fatalf(`
		types are different:
		expected: %T
		got: %T
		`,
			expected, actual,
		)

	}
}

func (a *assertion) NoError(err error) {
	if err != nil {
		a.t.Fatalf(`
		no error expected but got:
			%s
		`,
			err,
		)
	}
}
func Assert(t *testing.T) *assertion {
	return &assertion{t: t}
}
