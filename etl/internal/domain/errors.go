package domain

import "errors"

var (
	ErrConfigFileLoad = errors.New("error loading .env file")
	ErrConfigParse    = errors.New("error parsing config")
	ErrUnknownSource  = errors.New("unknown source")
)
