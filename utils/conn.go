package utils

import (
	"io"
	"net"
	"os"
)

type ErrorUnwrapper interface {
	Unwrap() error
}

func IsConnEOL(err error) bool {
	if err == nil {
		return false
	}
	if err == io.EOF || err == io.ErrClosedPipe || err == io.ErrUnexpectedEOF || err == net.ErrClosed || err == os.ErrClosed || err == os.ErrDeadlineExceeded {
		return true
	}
	if opError, ok := err.(*net.OpError); ok {
		return IsConnEOL(opError.Err)
	}
	if wrappedError, ok := err.(ErrorUnwrapper); ok {
		return IsConnEOL(wrappedError.Unwrap())
	}
	return false
}
