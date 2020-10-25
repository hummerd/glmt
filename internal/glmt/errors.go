package glmt

import "errors"

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
