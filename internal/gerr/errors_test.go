package gerr

import (
	"errors"
	"io"
	"net"
	"testing"
)

func TestNestedErrorIs(t *testing.T) {
	n := NewNestedError(net.ErrWriteToConnected, io.EOF)

	if !errors.Is(n, NestedError{}) {
		t.Fatal("error is NestedError")
	}

	if !errors.Is(n, net.ErrWriteToConnected) {
		t.Fatal("error is not ErrNotification")
	}

	if !errors.Is(n, io.EOF) {
		t.Fatal("error is not io.EOF")
	}
}

func TestNestedErrorAs(t *testing.T) {
	n := NewNestedError(&net.ParseError{}, &net.OpError{})

	var ne NestedError
	if !errors.As(n, &ne) {
		t.Fatal("error is not NestedError")
	}

	pe := &net.ParseError{}
	if !errors.As(n, &pe) {
		t.Fatal("error is not ErrNotification")
	}

	ope := &net.OpError{}
	if !errors.As(n, &ope) {
		t.Fatal("error is not io.EOF")
	}
}
