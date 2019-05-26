package common

import (
	"errors"
)

var (
	ErrOK   = errors.New("OK")
	ErrFail = errors.New("Fail")
)
