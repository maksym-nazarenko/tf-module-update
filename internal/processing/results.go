package processing

import (
	"fmt"
	"log"
	"strings"

	"github.com/maxim-nazarenko/tf-module-update/internal/processing/logging"
)

// Result wraps log message with a given log level
type Result struct {
	Message string
	Level   logging.Level
}

func (p *Result) String() string {
	return p.Message
}

// Results holds a set of Result messages alongside with errors
type Results struct {
	errors  []error
	level   logging.Level
	results []Result
}

// Append adds more items to the results set which might be rendered later
func (p *Results) Append(new ...interface{}) {
	for _, v := range new {
		if v == nil {
			continue
		}

		switch t := v.(type) {
		case error:
			// errors are handled in different way, so we can easily implement HasErrors()
			// and print all errors at the end to make this information more visible
			p.errors = append(p.errors, t)
		case Result:
			p.results = append(p.results, t)
		case *Result:
			p.results = append(p.results, *t)
		case Results:
			p.results = append(p.results, t.results...)
		case *Results:
			p.results = append(p.results, t.results...)
		default:
			log.Fatalf("unsupported result type: %T", t)
		}
	}
}

// String renders result records as string
func (p *Results) String() string {
	lines := make([]string, 0)
	for _, v := range p.results {
		if v.Level < p.level {
			continue
		}

		lines = append(lines, v.String())
	}
	for _, v := range p.errors {
		lines = append(lines, v.Error())
	}

	return strings.Join(lines, "\n")
}

// HasErrors indicates that there was an error
func (p *Results) HasErrors() bool {
	return len(p.errors) > 0
}

// LevelFromString converts string representation of log level to its typed version
func LevelFromString(logLevel string) (logging.Level, error) {
	level, ok := map[string]logging.Level{
		"trace": logging.TRACE,
		"debug": logging.DEBUG,
		"info":  logging.INFO,
		"warn":  logging.WARN,
		"error": logging.ERROR,
	}[strings.ToLower(logLevel)]

	if !ok {
		return logging.INFO, fmt.Errorf("unknown log level: %s", logLevel)
	}

	return level, nil
}

func NewResults(level logging.Level) *Results {
	return &Results{level: level}
}
