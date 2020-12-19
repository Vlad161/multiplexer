package handler

import (
	"errors"
)

var (
	ErrDecodeBody  = errors.New("can't decode body")
	ErrEmptyUrls   = errors.New("empty urls")
	ErrTooMuchUrls = errors.New("too much urls")
)
