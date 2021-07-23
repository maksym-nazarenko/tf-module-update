package processing

import (
	"errors"
	"testing"

	"github.com/maxim-nazarenko/tf-module-update/internal/processing/logging"
	"github.com/maxim-nazarenko/tf-module-update/internal/testhelpers"
)

func TestResultsString(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult string
		items          []interface{}
		loglevel       logging.Level
	}{
		{
			name: "info level, regular messages info level",
			expectedResult: `message 1
message 2`,
			items:    []interface{}{Result{Message: "message 1", Level: logging.INFO}, Result{Message: "message 2", Level: logging.INFO}},
			loglevel: logging.INFO,
		},
		{
			name: "info level, errors go to the bottom",
			expectedResult: `message 1
message 2
error 1
error 2`,
			items:    []interface{}{errors.New("error 1"), Result{Message: "message 1", Level: logging.INFO}, errors.New("error 2"), Result{Message: "message 2", Level: logging.INFO}},
			loglevel: logging.INFO,
		},
		{
			name: "info level, results level is warning",
			expectedResult: `message 2
error 1
error 2`,
			items:    []interface{}{errors.New("error 1"), Result{Message: "message 1", Level: logging.INFO}, errors.New("error 2"), Result{Message: "message 2", Level: logging.WARN}},
			loglevel: logging.WARN,
		},
	}
	assert := testhelpers.Assert(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results := NewResults(tc.loglevel)
			results.Append(tc.items...)
			assert.Equal(tc.expectedResult, results.String())
		})
	}
}

func TestResultsHasErrors(t *testing.T) {
	testCases := []struct {
		name           string
		items          []interface{}
		expectedResult bool
	}{
		{
			name:           "empty struct",
			expectedResult: false,
			items:          nil,
		},
		{
			name:           "info message",
			expectedResult: false,
			items:          []interface{}{Result{Message: "a regular message"}},
		},
		{
			name:           "error is added",
			expectedResult: true,
			items:          []interface{}{errors.New("regular error")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := testhelpers.Assert(t)
			r := Results{}
			r.Append(tc.items...)
			result := r.HasErrors()
			assert.Equal(tc.expectedResult, result)
		})
	}
}

func TestResultsLevelFromString(t *testing.T) {
	testCases := []struct {
		name           string
		stringLevel    string
		expectedResult logging.Level
		expectedError  error
	}{
		{
			name:           "lowercase debug",
			expectedResult: logging.DEBUG,
			stringLevel:    "debug",
		},
		{
			name:           "mixed case eRRor",
			expectedResult: logging.ERROR,
			stringLevel:    "eRRor",
		},
		{
			name:          "unknown level 'unknown'",
			expectedError: errors.New(""),
			stringLevel:   "unknown",
		},
	}
	assert := testhelpers.Assert(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			level, err := LevelFromString(tc.stringLevel)
			if tc.expectedError != nil {
				assert.SameType(tc.expectedError, err)
				return
			}
			assert.Equal(tc.expectedResult, level)
		})
	}
}
