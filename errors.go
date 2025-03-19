package muxedsocket

import "errors"

var (
	ErrSchemeNotSupported    = errors.New("scheme not supported")
	ErrRedialNotSupported    = errors.New("redial not supported")
	ErrInvalidChainingResult = errors.New("invalid chaining result")
	ErrOpNotSupported        = errors.New("not supported")
	// ErrIncompatibleChainOfLayers = errors.New("incompatible mix of layers")
)

type ErrMissingPart string

func (e ErrMissingPart) Error() string {
	return "missing part: " + string(e)
}

var _ error = ErrMissingPart("")
