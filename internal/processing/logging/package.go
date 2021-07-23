package logging

// Level represents logging level
type Level int

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
)
