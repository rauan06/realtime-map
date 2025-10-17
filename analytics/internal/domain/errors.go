package domain

import "errors"

var (
	ErrConfigFileLode   = errors.New("error loading .env file")
	ErrConfigError      = errors.New("loading config")
	ErrInvalidByteOrder = errors.New("invalid byte order")
	ErrInvalidType      = errors.New("invalid type")
)
