package processing

import "github.com/maxim-nazarenko/tf-module-update/internal/processing/logging"

// ResultFactory wraps building of Result messages with convenient logging functions
type ResultFactory struct {
}

func (rf *ResultFactory) Trace(msg string) Result {
	return Result{Message: msg, Level: logging.TRACE}
}

func (rf *ResultFactory) Debug(msg string) Result {
	return Result{Message: msg, Level: logging.DEBUG}
}

func (rf *ResultFactory) Info(msg string) Result {
	return Result{Message: msg, Level: logging.INFO}
}

func (rf *ResultFactory) Warn(msg string) Result {
	return Result{Message: msg, Level: logging.WARN}
}

func (rf *ResultFactory) Error(msg string) Result {
	return Result{Message: msg, Level: logging.ERROR}
}

func NewResultFactory() *ResultFactory {
	return &ResultFactory{}
}
