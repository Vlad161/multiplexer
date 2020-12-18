package handler

import (
	"errors"
)

var (
	ErrDecodeBody  = errors.New("can't decode body")
	ErrTooMuchUrls = errors.New("too much urls")
)
