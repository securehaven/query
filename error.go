package query

import "errors"

type ParsingError struct {
	Errors []error
}

func (e ParsingError) Error() string {
	return errors.Join(e.Errors...).Error()
}
