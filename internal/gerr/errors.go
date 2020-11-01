// Package gerr contains special errors
package gerr

import (
	"errors"
	"strings"
)

func NewNestedError(wrap, cause error) error {
	if wrap == nil || cause == nil {
		panic("wrap or cause can not be nil")
	}

	return NestedError{
		wrap:  wrap,
		cause: cause,
	}
}

// NestedError is usefull when you want to wrap unknown error
// with predefined error
type NestedError struct {
	wrap  error
	cause error
}

func (ne NestedError) Error() string {
	return ne.wrap.Error() + ": " + ne.cause.Error()
}

func (ne NestedError) Unwrap() error {
	return ne.wrap
}

func (ne NestedError) Is(err error) bool {
	switch err.(type) {
	case NestedError:
		return true
	}

	return errors.Is(ne.wrap, err) ||
		errors.Is(ne.cause, err)
}

func (ne NestedError) As(as interface{}) bool {
	return errors.As(ne.wrap, as) ||
		errors.As(ne.cause, as)
}

func NewMultiError(errs ...error) error {
	var actualErrs []error
	for _, err := range errs {
		if err == nil {
			continue
		}

		actualErrs = append(actualErrs, err)
	}

	if len(actualErrs) == 0 {
		return nil
	}

	return MultiErr{
		Errs: actualErrs,
	}
}

type MultiErr struct {
	Errs []error
}

func (me MultiErr) Unwrap() error {
	if len(me.Errs) == 0 {
		return nil
	}

	return me.Errs[0]
}

func (me MultiErr) Is(err error) bool {
	switch err.(type) {
	case MultiErr:
		return true
	}

	for _, e := range me.Errs {
		if errors.Is(e, err) {
			return true
		}
	}

	return false
}

func (me MultiErr) As(as interface{}) bool {
	for _, e := range me.Errs {
		if errors.As(e, as) {
			return true
		}
	}

	return false
}

func (me MultiErr) Error() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString("errors: ")

	for i, e := range me.Errs {
		sb.WriteString(e.Error())
		if i < len(me.Errs)-1 {
			sb.WriteString("; ")
		}
	}

	return sb.String()
}
