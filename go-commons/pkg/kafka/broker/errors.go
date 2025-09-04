package broker

import "errors"

var (
	// ErrTimeout -.
	ErrTimeout = errors.New("timeout")
	// ErrInternalServer -.
	ErrInternalServer = errors.New("internal server error")
	// ErrBadHandler -.
	ErrBadTopc = errors.New("unrecognized topic")
)

// Success -.
const Success = "success"
