package errors

import (
	"errors"
)

type ErrorList []error

func NewErrors() *ErrorList {
	return &ErrorList{}
}

func (e *ErrorList) Join(err string) error {
	if len(*e) == 0 {
		return nil
	}

	if e.Length() == 1 {
		return (*e)[0]
	}

	result := make([]error, 1, 1+e.Length())
	result[0] = errors.New(err)
	result = append(result, *e...)
	return errors.Join(result...)
}

func (e *ErrorList) Add(err error) {
	e.Append(err)
}

func (e *ErrorList) Append(err ...error) {
	if len(err) == 0 {
		return
	} else if len(err) == 1 && err[0] == nil {
		return
	}

	*e = append(*e, err...)
}

func (e *ErrorList) Errors() []error {
	return *e
}

func (e *ErrorList) Length() int {
	return len(*e)
}
