package tcpnetwork

import (
	"errors"
)

var (
	kErrConnIsClosed    = errors.New("Connection is closed")
	kErrConnSendTimeout = errors.New("Connection send timeout")
)
