package module

type InvalidSourceFormatError struct {
	text string
}

func (e *InvalidSourceFormatError) Error() string {
	return e.text
}
